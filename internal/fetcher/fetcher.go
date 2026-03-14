package fetcher

import (
	"context"
	"fmt"
	"log"
	"maps"
	"math/rand"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/time/rate"

	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source"
)

const (
	defaultFetchTimeout        = 30 * time.Second
	defaultRedditHostSpacing   = 3 * time.Second
	defaultRedditStartupJitter = 10 * time.Second
	defaultRedditBackoffCap    = 2 * time.Hour
	defaultSourceInterval      = 30 * time.Minute
)

type Fetcher struct {
	sources        []sourceWithConfig
	feed           *models.Feed
	mu             sync.RWMutex
	subscribers    map[chan struct{}]struct{}
	subMu          sync.Mutex
	maxSubscribers int
	limiters       map[string]*rate.Limiter
	limiterMu      sync.Mutex
	hostLimiters   map[string]*rate.Limiter
	hostLimiterMu  sync.Mutex
	minInterval    time.Duration
	wg             sync.WaitGroup
	closed         bool
	closedMu       sync.Mutex
	subOnce        map[chan struct{}]*sync.Once
	rng            *rand.Rand
	rngMu          sync.Mutex
}

// Config holds configuration for the fetcher.
type Config struct {
	MaxSubscribers   int
	MinFetchInterval time.Duration
}

type sourceWithConfig struct {
	source            source.Source
	interval          time.Duration
	intervalJitter    time.Duration
	host              string
	isReddit          bool
	startupStaggerMax time.Duration
	failureBackoffCap time.Duration
}

// New creates a new Fetcher with default config.
func New(sources []sourceWithConfig) *Fetcher {
	return NewWithConfig(sources, Config{
		MaxSubscribers:   1000,
		MinFetchInterval: 0,
	})
}

func newFeed() *models.Feed {
	return &models.Feed{
		Items:        make([]models.Item, 0),
		UpdatedAt:    time.Now(),
		Errors:       make(map[string]string),
		SourceStates: make(map[string]models.SourceState),
	}
}

func NewWithConfig(sources []sourceWithConfig, cfg Config) *Fetcher {
	return &Fetcher{
		sources: sources,
		feed:    newFeed(),

		subscribers:    make(map[chan struct{}]struct{}),
		subOnce:        make(map[chan struct{}]*sync.Once),
		maxSubscribers: cfg.MaxSubscribers,
		limiters:       make(map[string]*rate.Limiter),
		hostLimiters:   make(map[string]*rate.Limiter),
		minInterval:    cfg.MinFetchInterval,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewFromConfigs creates a new Fetcher from source configs.
func NewFromConfigs(configs []models.SourceConfig, minFetchInterval int, maxSubscribers int) *Fetcher {
	sources := make([]sourceWithConfig, 0, len(configs))

	for _, cfg := range configs {
		var src source.Source
		switch cfg.Type {
		case "rss":
			src = source.NewFeedSource(cfg.Name, cfg.URL, "rss", false, false, false)
		case "hnalgolia":
			src = source.NewHNAlgoliaSource(cfg.Name, cfg.URL, cfg.Type)
		case "reddit":
			src = source.NewFeedSource(cfg.Name, cfg.URL, "reddit", false, cfg.IgnoreDays, cfg.IsChronological)
		case "lobsters":
			src = source.NewFeedSource(cfg.Name, cfg.URL, "lobsters", true, false, false)
		case "4plebs":
			src = source.NewFourPlebsSource(cfg.Name, cfg.URL, cfg.Limit, cfg.NSFW)
		case "desuarchive":
			src = source.NewDesuArchiveSource(cfg.Name, cfg.URL, cfg.Limit, cfg.NSFW)
		case "wikipedia":
			src = source.NewWikipediaSource(cfg.Name, cfg.URL, cfg.Limit)
		default:
			log.Printf("Unknown source type: %s", cfg.Type)
			continue
		}

		sources = append(sources, sourceWithConfig{
			source:            src,
			interval:          time.Duration(cfg.Interval) * time.Second,
			intervalJitter:    time.Duration(cfg.IntervalJitter) * time.Second,
			host:              sourceHost(cfg.URL),
			isReddit:          cfg.Type == "reddit",
			startupStaggerMax: defaultRedditStartupJitter,
			failureBackoffCap: defaultRedditBackoffCap,
		})
	}

	return &Fetcher{
		sources: sources,
		feed:    newFeed(),

		subscribers:    make(map[chan struct{}]struct{}),
		subOnce:        make(map[chan struct{}]*sync.Once),
		maxSubscribers: maxSubscribers,
		limiters:       make(map[string]*rate.Limiter),
		hostLimiters:   make(map[string]*rate.Limiter),
		minInterval:    time.Duration(minFetchInterval) * time.Second,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins fetching from all sources in background.
func (f *Fetcher) Start(ctx context.Context) {
	f.initSourceStates()

	for _, sc := range f.sources {
		f.wg.Add(1)
		go func(sc sourceWithConfig) {
			defer f.wg.Done()
			f.fetchLoop(ctx, sc)
		}(sc)
	}
}

// fetchLoop runs a fetch loop for a single source with semi-random intervals.
func (f *Fetcher) fetchLoop(ctx context.Context, sc sourceWithConfig) {
	if ctx.Err() != nil {
		return
	}

	if !f.waitInitialDelay(ctx, sc) {
		return
	}

	for {
		if ctx.Err() != nil {
			return
		}

		f.fetchSource(ctx, sc)
		f.notifySubscribers()

		delay := f.nextDelay(sc)
		f.logNextFetch(sc, delay)

		timer := time.NewTimer(delay)
		select {
		case <-ctx.Done():
			timer.Stop()
			return
		case <-timer.C:
		}
	}
}

func (f *Fetcher) waitInitialDelay(ctx context.Context, sc sourceWithConfig) bool {
	if !sc.isReddit || sc.startupStaggerMax <= 0 {
		return true
	}

	delay := f.randomDuration(sc.startupStaggerMax)
	if delay <= 0 {
		return true
	}

	log.Printf("Initial stagger for %s: first fetch %s", sc.source.Name(), delay.Round(time.Second))

	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

func (f *Fetcher) nextDelay(sc sourceWithConfig) time.Duration {
	base := sc.interval
	if base <= 0 {
		base = defaultSourceInterval
	}
	if sc.intervalJitter > 0 {
		base += f.randomDuration(sc.intervalJitter)
	}
	if !sc.isReddit {
		return base
	}

	failures := f.consecutiveFailures(sc.source.Name())
	if failures <= 1 {
		return base
	}

	multiplier := 1
	for i := 1; i < failures; i++ {
		if multiplier >= 64 {
			break
		}
		multiplier *= 2
	}

	delay := time.Duration(multiplier) * base
	if sc.failureBackoffCap > 0 && delay > sc.failureBackoffCap {
		return sc.failureBackoffCap
	}
	return delay
}

func (f *Fetcher) fetchSource(ctx context.Context, sc sourceWithConfig) {
	src := sc.source

	if f.minInterval > 0 {
		limiter := f.getLimiter(src)
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Rate limiting %s: need to wait", src.Name())
			return
		}
	}

	if limiter := f.getHostLimiter(sc); limiter != nil {
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Host limiting %s (%s): %v", src.Name(), sc.host, err)
			return
		}
	}

	start := time.Now()
	attemptAt := start
	f.markAttempt(sc, attemptAt)
	log.Printf("Fetching from %s (%s)", src.Name(), src.Type())

	fetchCtx, fetchCancel := context.WithTimeout(ctx, defaultFetchTimeout)
	defer fetchCancel()

	items, err := src.Fetch(fetchCtx)
	duration := time.Since(start)

	if err != nil {
		failures := f.markFailure(sc, attemptAt, err)
		log.Printf("Error fetching from %s (host=%s, duration=%s, failures=%d): %v", src.Name(), sc.host, duration.Round(time.Millisecond), failures, err)
		return
	}

	f.markSuccess(sc, attemptAt, items)
	log.Printf("Fetched %d items from %s (host=%s, duration=%s)", len(items), src.Name(), sc.host, duration.Round(time.Millisecond))
}

func (f *Fetcher) getLimiter(src source.Source) *rate.Limiter {
	f.limiterMu.Lock()
	defer f.limiterMu.Unlock()

	if limiter, exists := f.limiters[src.Name()]; exists {
		return limiter
	}

	r := rate.Every(f.minInterval)
	limiter := rate.NewLimiter(r, 2)
	f.limiters[src.Name()] = limiter
	return limiter
}

func (f *Fetcher) getHostLimiter(sc sourceWithConfig) *rate.Limiter {
	if !sc.isReddit || sc.host == "" {
		return nil
	}

	f.hostLimiterMu.Lock()
	defer f.hostLimiterMu.Unlock()

	if limiter, exists := f.hostLimiters[sc.host]; exists {
		return limiter
	}

	limiter := rate.NewLimiter(rate.Every(defaultRedditHostSpacing), 1)
	f.hostLimiters[sc.host] = limiter
	return limiter
}

func (f *Fetcher) initSourceStates() {
	f.mu.Lock()
	defer f.mu.Unlock()

	for _, sc := range f.sources {
		state := f.ensureSourceStateLocked(sc)
		f.feed.SourceStates[sc.source.Name()] = state
	}
}

func (f *Fetcher) ensureSourceStateLocked(sc sourceWithConfig) models.SourceState {
	if state, exists := f.feed.SourceStates[sc.source.Name()]; exists {
		if state.Type == "" {
			state.Type = sc.source.Type()
		}
		if state.Host == "" {
			state.Host = sc.host
		}
		if state.Name == "" {
			state.Name = sc.source.Name()
		}
		return state
	}

	return models.SourceState{
		Name:  sc.source.Name(),
		Type:  sc.source.Type(),
		Host:  sc.host,
		Stale: false,
	}
}

func (f *Fetcher) markAttempt(sc sourceWithConfig, at time.Time) {
	f.mu.Lock()
	defer f.mu.Unlock()

	state := f.ensureSourceStateLocked(sc)
	state.LastAttemptAt = at
	f.feed.SourceStates[sc.source.Name()] = state
}

func (f *Fetcher) markFailure(sc sourceWithConfig, at time.Time, err error) int {
	f.mu.Lock()
	defer f.mu.Unlock()

	state := f.ensureSourceStateLocked(sc)
	state.LastAttemptAt = at
	state.LastError = err.Error()
	state.ConsecutiveFailures++
	state.Stale = true
	f.feed.SourceStates[sc.source.Name()] = state
	f.feed.Errors[sc.source.Name()] = err.Error()
	f.feed.UpdatedAt = time.Now()

	return state.ConsecutiveFailures
}

func (f *Fetcher) markSuccess(sc sourceWithConfig, at time.Time, items []models.Item) {
	f.mu.Lock()
	defer f.mu.Unlock()

	newItems := make([]models.Item, 0, len(f.feed.Items)+len(items))
	for _, item := range f.feed.Items {
		if item.SourceName != sc.source.Name() {
			newItems = append(newItems, item)
		}
	}
	newItems = append(newItems, items...)
	f.feed.Items = newItems

	state := f.ensureSourceStateLocked(sc)
	state.LastAttemptAt = at
	state.LastSuccessAt = at
	state.LastError = ""
	state.ConsecutiveFailures = 0
	state.Stale = false
	f.feed.SourceStates[sc.source.Name()] = state

	delete(f.feed.Errors, sc.source.Name())
	f.feed.UpdatedAt = time.Now()
}

func (f *Fetcher) consecutiveFailures(sourceName string) int {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if state, exists := f.feed.SourceStates[sourceName]; exists {
		return state.ConsecutiveFailures
	}
	return 0
}

func (f *Fetcher) logNextFetch(sc sourceWithConfig, delay time.Duration) {
	failures := f.consecutiveFailures(sc.source.Name())
	backoff := sc.isReddit && failures > 1
	log.Printf("Next fetch for %s in %s (backoff=%t, failures=%d)", sc.source.Name(), delay.Round(time.Second), backoff, failures)
}

func (f *Fetcher) randomDuration(max time.Duration) time.Duration {
	if max <= 0 {
		return 0
	}

	f.rngMu.Lock()
	defer f.rngMu.Unlock()
	return time.Duration(f.rng.Int63n(int64(max)))
}

func sourceHost(rawURL string) string {
	parsed, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	return strings.ToLower(parsed.Hostname())
}

func (f *Fetcher) Shutdown() error {
	f.closedMu.Lock()
	f.closed = true
	f.closedMu.Unlock()

	f.wg.Wait()
	return nil
}

// Subscribe returns a channel that receives notifications when feed updates.
func (f *Fetcher) Subscribe() (chan struct{}, error) {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	if f.maxSubscribers > 0 && len(f.subscribers) >= f.maxSubscribers {
		return nil, fmt.Errorf("subscriber limit reached (max: %d)", f.maxSubscribers)
	}

	f.closedMu.Lock()
	isClosed := f.closed
	f.closedMu.Unlock()

	if isClosed {
		return nil, fmt.Errorf("fetcher is closed")
	}

	ch := make(chan struct{}, 1)
	f.subscribers[ch] = struct{}{}
	f.subOnce[ch] = &sync.Once{}
	return ch, nil
}

// Unsubscribe removes a subscriber channel.
func (f *Fetcher) Unsubscribe(ch chan struct{}) {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	if _, exists := f.subscribers[ch]; exists {
		delete(f.subscribers, ch)
		if once, ok := f.subOnce[ch]; ok {
			once.Do(func() {
				close(ch)
			})
		}
		delete(f.subOnce, ch)
	}
}

// notifySubscribers sends update notifications to all subscribers.
func (f *Fetcher) notifySubscribers() {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	for ch := range f.subscribers {
		select {
		case ch <- struct{}{}:
		default:
		}
	}
}

// GetFeed returns a copy of the current feed.
func (f *Fetcher) GetFeed() models.Feed {
	f.mu.RLock()
	defer f.mu.RUnlock()

	errorsCopy := make(map[string]string, len(f.feed.Errors))
	maps.Copy(errorsCopy, f.feed.Errors)

	statesCopy := make(map[string]models.SourceState, len(f.feed.SourceStates))
	maps.Copy(statesCopy, f.feed.SourceStates)

	return models.Feed{
		Items:        append([]models.Item(nil), f.feed.Items...),
		UpdatedAt:    f.feed.UpdatedAt,
		Errors:       errorsCopy,
		SourceStates: statesCopy,
	}
}
