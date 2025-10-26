package fetcher

import (
	"context"
	"log"
	"maps"
	"math/rand"
	"sync"
	"time"

	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source"
)

// Fetcher manages fetching from multiple sources with semi-random intervals
type Fetcher struct {
	sources     []sourceWithConfig
	feed        *models.Feed
	mu          sync.RWMutex
	subscribers []chan struct{}
	subMu       sync.Mutex
}

type sourceWithConfig struct {
	source         source.Source
	interval       time.Duration
	intervalJitter time.Duration
}

// New creates a new Fetcher
func New(sources []sourceWithConfig) *Fetcher {
	return &Fetcher{
		sources: sources,
		feed: &models.Feed{
			Items:     make([]models.Item, 0),
			UpdatedAt: time.Now(),
			Errors:    make(map[string]string),
		},
	}
}

// NewFromConfigs creates a new Fetcher from source configs
func NewFromConfigs(configs []models.SourceConfig) *Fetcher {
	sources := make([]sourceWithConfig, 0, len(configs))

	for _, cfg := range configs {
		var src source.Source
		switch cfg.Type {
		case "rss":
			src = source.NewRSSSource(cfg.Name, cfg.URL)
		case "hnrss":
			src = source.NewHNRSSSource(cfg.Name, cfg.URL)
		case "reddit":
			src = source.NewRedditSource(cfg.Name, cfg.URL, cfg.IgnoreDays)
		case "lobsters":
			src = source.NewLobstersSource(cfg.Name, cfg.URL)
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

	return New(sources)
}

// Start begins fetching from all sources
func (f *Fetcher) Start(ctx context.Context) {
	var wg sync.WaitGroup

	for _, sc := range f.sources {
		wg.Add(1)
		go func(sc sourceWithConfig) {
			defer wg.Done()
			f.fetchLoop(ctx, sc)
		}(sc)
	}

	wg.Wait()
}

// fetchLoop runs a fetch loop for a single source with semi-random intervals
func (f *Fetcher) fetchLoop(ctx context.Context, sc sourceWithConfig) {
	// Fetch immediately on start
	f.fetchSource(ctx, sc.source)

	ticker := time.NewTicker(f.nextInterval(sc))
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
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
	jitter := time.Duration(rand.Int63n(int64(sc.intervalJitter)))
	return sc.interval + jitter
}

// fetchSource fetches from a single source and updates the feed
func (f *Fetcher) fetchSource(ctx context.Context, src source.Source) {
	log.Printf("Fetching from %s (%s)", src.Name(), src.Type())

	items, err := src.Fetch(ctx)

	f.mu.Lock()
	defer f.mu.Unlock()

	if err != nil {
		log.Printf("Error fetching from %s: %v", src.Name(), err)
		// Store error so frontend can display it
		f.feed.Errors[src.Name()] = err.Error()
		f.feed.UpdatedAt = time.Now()
		f.notifySubscribers()
		return
	}

	// Clear any previous error for this source
	delete(f.feed.Errors, src.Name())

	// Remove old items from this source
	newItems := make([]models.Item, 0, len(f.feed.Items))
	for _, item := range f.feed.Items {
		if item.SourceName != src.Name() {
			newItems = append(newItems, item)
		}
	}

	// Add new items
	newItems = append(newItems, items...)
	f.feed.Items = newItems
	f.feed.UpdatedAt = time.Now()

	log.Printf("Fetched %d items from %s", len(items), src.Name())

	// Notify all subscribers
	f.notifySubscribers()
}

// Subscribe returns a channel that receives notifications when feed updates
func (f *Fetcher) Subscribe() chan struct{} {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	ch := make(chan struct{}, 1)
	f.subscribers = append(f.subscribers, ch)
	return ch
}

// Unsubscribe removes a subscriber channel
func (f *Fetcher) Unsubscribe(ch chan struct{}) {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	for i, sub := range f.subscribers {
		if sub == ch {
			f.subscribers = append(f.subscribers[:i], f.subscribers[i+1:]...)
			close(ch)
			break
		}
	}
}

// notifySubscribers sends update notifications to all subscribers
func (f *Fetcher) notifySubscribers() {
	f.subMu.Lock()
	defer f.subMu.Unlock()

	for _, ch := range f.subscribers {
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
