package aggregator

import (
	"sort"

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

// LimitPerSource limits the number of items per source to the given limit.
func (a *Aggregate) LimitPerSource(limit int) *Aggregate {
	if limit <= 0 {
		return a
	}

	grouped := a.GroupBySource()
	filtered := make([]models.Item, 0, len(a.Items))
	seen := make(map[string]bool)

	for _, item := range a.Items {
		if seen[item.SourceName] {
			continue
		}
		seen[item.SourceName] = true

		items := grouped[item.SourceName]
		if len(items) > limit {
			items = items[:limit]
		}
		filtered = append(filtered, items...)
	}

	return &Aggregate{
		Items: filtered,
	}
}
