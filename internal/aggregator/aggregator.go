package aggregator

import (
	"sort"
	"time"

	"github.com/ppowo/feedlet/internal/models"
)

// Aggregate processes and sorts items from a feed
type Aggregate struct {
	Items []models.Item
}

// Process takes a feed and returns an aggregated, sorted view
func Process(feed models.Feed) *Aggregate {
	// Sort items by published date (newest first)
	items := append([]models.Item(nil), feed.Items...)
	sort.Slice(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})

	return &Aggregate{
		Items: items,
	}
}

// GroupBySource groups items by their source name
func (a *Aggregate) GroupBySource() map[string][]models.Item {
	grouped := make(map[string][]models.Item)
	for _, item := range a.Items {
		grouped[item.SourceName] = append(grouped[item.SourceName], item)
	}
	return grouped
}

// GroupByDate groups items by date (YYYY-MM-DD)
func (a *Aggregate) GroupByDate() map[string][]models.Item {
	grouped := make(map[string][]models.Item)
	for _, item := range a.Items {
		date := item.Published.Format("2006-01-02")
		grouped[date] = append(grouped[date], item)
	}
	return grouped
}

// Limit returns only the first n items
func (a *Aggregate) Limit(n int) *Aggregate {
	if n > len(a.Items) {
		n = len(a.Items)
	}
	return &Aggregate{
		Items: a.Items[:n],
	}
}

// FilterByAge returns only items published within the last n days
// Items with IgnoreDays=true are always included
func (a *Aggregate) FilterByAge(days int) *Aggregate {
	cutoff := time.Now().AddDate(0, 0, -days)
	filtered := make([]models.Item, 0, len(a.Items))

	for _, item := range a.Items {
		if item.IgnoreDays || item.Published.After(cutoff) {
			filtered = append(filtered, item)
		}
	}

	return &Aggregate{
		Items: filtered,
	}
}

// FilterBySourceDays filters items based on per-source day limits
// Items with IgnoreDays=true are always included
func (a *Aggregate) FilterBySourceDays(sourceDays map[string]int) *Aggregate {
	filtered := make([]models.Item, 0, len(a.Items))

	for _, item := range a.Items {
		if item.IgnoreDays {
			// Always include items that ignore days
			filtered = append(filtered, item)
			continue
		}

		days, hasDays := sourceDays[item.SourceName]
		if !hasDays || days <= 0 {
			// No day limit for this source, include all
			filtered = append(filtered, item)
			continue
		}

		cutoff := time.Now().AddDate(0, 0, -days)
		if item.Published.After(cutoff) {
			filtered = append(filtered, item)
		}
	}

	return &Aggregate{
		Items: filtered,
	}
}

// LimitPerSource limits the number of items per source
// If limit is 0 or negative, no limiting is applied
func (a *Aggregate) LimitPerSource(limits map[string]int) *Aggregate {
	if len(limits) == 0 {
		return a
	}

	// Track count per source
	counts := make(map[string]int)
	filtered := make([]models.Item, 0, len(a.Items))

	for _, item := range a.Items {
		limit, hasLimit := limits[item.SourceName]
		if !hasLimit || limit <= 0 {
			// No limit for this source
			filtered = append(filtered, item)
			continue
		}

		if counts[item.SourceName] < limit {
			filtered = append(filtered, item)
			counts[item.SourceName]++
		}
	}

	return &Aggregate{
		Items: filtered,
	}
}
