package aggregator

import (
	"sort"
	"time"

	"github.com/ppowo/feedlet/internal/models"
)

// Aggregate processes items from a feed.
type Aggregate struct {
	Items []models.Item
}

// Process takes a feed and returns an aggregate view.
func Process(feed models.Feed) *Aggregate {
	return &Aggregate{
		Items: append([]models.Item(nil), feed.Items...),
	}
}

func sortByPublished(items []models.Item) {
	if len(items) < 2 {
		return
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})
}

// GroupBySource groups items by their source name.
func (a *Aggregate) GroupBySource() map[string][]models.Item {
	grouped := make(map[string][]models.Item)

	for _, item := range a.Items {
		grouped[item.SourceName] = append(grouped[item.SourceName], item)
	}

	for _, items := range grouped {
		sortByPublished(items)
	}

	return grouped
}

// GroupByDate groups items by date (YYYY-MM-DD).
func (a *Aggregate) GroupByDate() map[string][]models.Item {
	grouped := make(map[string][]models.Item)
	for _, item := range a.Items {
		date := item.Published.Format("2006-01-02")
		grouped[date] = append(grouped[date], item)
	}

	for _, items := range grouped {
		sortByPublished(items)
	}

	return grouped
}

// Limit returns only the first n items.
func (a *Aggregate) Limit(n int) *Aggregate {
	if n <= 0 {
		return &Aggregate{Items: []models.Item{}}
	}
	if n > len(a.Items) {
		n = len(a.Items)
	}
	return &Aggregate{
		Items: a.Items[:n],
	}
}

// FilterByAge returns only items published within the last n days.
// Items with IgnoreDays=true are always included.
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

// FilterBySourceDays filters items based on per-source day limits.
// Items with IgnoreDays=true are always included.
func (a *Aggregate) FilterBySourceDays(sourceDays map[string]int) *Aggregate {
	filtered := make([]models.Item, 0, len(a.Items))

	for _, item := range a.Items {
		if item.IgnoreDays {
			filtered = append(filtered, item)
			continue
		}

		days, hasDays := sourceDays[item.SourceName]
		if !hasDays || days <= 0 {
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

// LimitPerSource limits the number of items per source.
// If limit is 0 or negative, no limiting is applied.
func (a *Aggregate) LimitPerSource(limits map[string]int) *Aggregate {
	if len(limits) == 0 {
		return a
	}

	grouped := a.GroupBySource()
	filtered := make([]models.Item, 0, len(a.Items))
	seen := make(map[string]struct{})

	for _, item := range a.Items {
		sourceName := item.SourceName
		if _, ok := seen[sourceName]; ok {
			continue
		}
		seen[sourceName] = struct{}{}

		items := grouped[sourceName]
		if limit := limits[sourceName]; limit > 0 && limit < len(items) {
			items = items[:limit]
		}
		filtered = append(filtered, items...)
	}

	return &Aggregate{
		Items: filtered,
	}
}
