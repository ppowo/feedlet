package fetcher

import (
	"context"
	"fmt"
	"golang.org/x/time/rate"
	"log"
	"maps"
	"math/rand"
	"sync"
	"time"

	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source"
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
	minInterval    time.Duration
	wg             sync.WaitGroup
	closed         bool
	closedMu       sync.Mutex
	subOnce        map[chan struct{}]*sync.Once
	rng            *rand.Rand
	rngMu          sync.Mutex
}

// Config holds configuration for the fetcher
type Config struct {
	MaxSubscribers   int
	MinFetchInterval time.Duration
}

type sourceWithConfig struct {
	source         source.Source
	interval       time.Duration
	intervalJitter time.Duration
}

// New creates a new Fetcher with default config
func New(sources []sourceWithConfig) *Fetcher {
	return NewWithConfig(sources, Config{
		MaxSubscribers:   1000,
		MinFetchInterval: 0,
	})
}

func NewWithConfig(sources []sourceWithConfig, cfg Config) *Fetcher {
	return &Fetcher{
		sources: sources,
		feed: &models.Feed{
			Items:     make([]models.Item, 0),
			UpdatedAt: time.Now(),
			Errors:    make(map[string]string),
		},
		subscribers:    make(map[chan struct{}]struct{}),
		subOnce:        make(map[chan struct{}]*sync.Once),
		maxSubscribers: cfg.MaxSubscribers,
		limiters:       make(map[string]*rate.Limiter),
		minInterval:    cfg.MinFetchInterval,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// NewFromConfigs creates a new Fetcher from source configs
func NewFromConfigs(configs []models.SourceConfig, minFetchInterval int, maxSubscribers int) *Fetcher {
	sources := make([]sourceWithConfig, 0, len(configs))

	for _, cfg := range configs {
		var src source.Source
		switch cfg.Type {
		case "rss":
			src = source.NewFeedSource(cfg.Name, cfg.URL, "rss", false, false, false)
		case "hnrss":
			src = source.NewFeedSource(cfg.Name, cfg.URL, "hnrss", true, false, false)
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
			source:         src,
			interval:       time.Duration(cfg.Interval) * time.Second,
			intervalJitter: time.Duration(cfg.IntervalJitter) * time.Second,
		})
	}

	return &Fetcher{
		sources: sources,
		feed: &models.Feed{
			Items:     make([]models.Item, 0),
			UpdatedAt: time.Now(),
			Errors:    make(map[string]string),
		},
		subscribers:    make(map[chan struct{}]struct{}),
		subOnce:        make(map[chan struct{}]*sync.Once),
		maxSubscribers: maxSubscribers,
		limiters:       make(map[string]*rate.Limiter),
		minInterval:    time.Duration(minFetchInterval) * time.Second,
		rng:            rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Start begins fetching from all sources in background
func (f *Fetcher) Start(ctx context.Context) {
	for _, sc := range f.sources {
		f.wg.Add(1)
		go func(sc sourceWithConfig) {
			defer f.wg.Done()
			f.fetchLoop(ctx, sc)
		}(sc)
	}
}

// fetchLoop runs a fetch loop for a single source with semi-random intervals
func (f *Fetcher) fetchLoop(ctx context.Context, sc sourceWithConfig) {
	// Check if context is already cancelled
	if ctx.Err() != nil {
		return
	}

	// Fetch immediately on start
	f.fetchSource(ctx, sc.source)

	ticker := time.NewTicker(f.nextInterval(sc))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Check context before fetching
			if ctx.Err() != nil {
				return
			}
			f.fetchSource(ctx, sc.source)
			ticker.Reset(f.nextInterval(sc))
		}
	}
}

// nextInterval calculates the next interval with jitter
func (f *Fetcher) nextInterval(sc sourceWithConfig) time.Duration {
	if sc.intervalJitter == 0 {
		return sc.interval
	}
	// Use shared random generator with mutex protection
	f.rngMu.Lock()
	defer f.rngMu.Unlock()
	jitter := time.Duration(f.rng.Int63n(int64(sc.intervalJitter)))
	return sc.interval + jitter
}

func (f *Fetcher) fetchSource(ctx context.Context, src source.Source) {
	log.Printf("Fetching from %s (%s)", src.Name(), src.Type())

	if f.minInterval > 0 {
		limiter := f.getLimiter(src)
		if err := limiter.Wait(ctx); err != nil {
			log.Printf("Rate limiting %s: need to wait", src.Name())
			return
		}
	}

	fetchCtx, fetchCancel := context.WithTimeout(ctx, 30*time.Second)
	defer fetchCancel()

	items, err := src.Fetch(fetchCtx)

	f.mu.Lock()
	defer f.mu.Unlock()

	if err != nil {
		log.Printf("Error fetching from %s: %v", src.Name(), err)
		f.feed.Errors[src.Name()] = err.Error()
		f.feed.UpdatedAt = time.Now()
		f.notifySubscribers()
		return
	}

	delete(f.feed.Errors, src.Name())

	newItems := make([]models.Item, 0, len(f.feed.Items))
	for _, item := range f.feed.Items {
		if item.SourceName != src.Name() {
			newItems = append(newItems, item)
		}
	}

	newItems = append(newItems, items...)
	f.feed.Items = newItems
	f.feed.UpdatedAt = time.Now()

	log.Printf("Fetched %d items from %s", len(items), src.Name())

	f.notifySubscribers()
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

func (f *Fetcher) Shutdown() error {
	f.closedMu.Lock()
	f.closed = true
	f.closedMu.Unlock()

	f.wg.Wait()
	return nil
}

// Subscribe returns a channel that receives notifications when feed updates
func (f *Fetcher) Subscribe() (chan struct{}, error) {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	// Check if we've reached the subscriber limit
	if f.maxSubscribers > 0 && len(f.subscribers) >= f.maxSubscribers {
		return nil, fmt.Errorf("subscriber limit reached (max: %d)", f.maxSubscribers)
	}

	// Check if fetcher is closed (under its own mutex)
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

// Unsubscribe removes a subscriber channel
func (f *Fetcher) Unsubscribe(ch chan struct{}) {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	if _, exists := f.subscribers[ch]; exists {
		delete(f.subscribers, ch)
		// Use sync.Once to prevent double-close
		if once, ok := f.subOnce[ch]; ok {
			once.Do(func() {
				close(ch)
			})
		}
		delete(f.subOnce, ch)
	}
}

// notifySubscribers sends update notifications to all subscribers
func (f *Fetcher) notifySubscribers() {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	for ch := range f.subscribers {
		select {
		case ch <- struct{}{}:
		default:
			// Channel full, skip
		}
	}
}

// GetFeed returns a copy of the current feed
func (f *Fetcher) GetFeed() models.Feed {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Copy errors map
	errorsCopy := make(map[string]string, len(f.feed.Errors))
	maps.Copy(errorsCopy, f.feed.Errors)

	// Return a copy to avoid race conditions
	return models.Feed{
		Items:     append([]models.Item(nil), f.feed.Items...),
		UpdatedAt: f.feed.UpdatedAt,
		Errors:    errorsCopy,
	}
}
