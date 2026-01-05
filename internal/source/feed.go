package source

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
)

// FeedSource implements the Source interface for RSS/Atom feeds
// Consolidates rss, reddit, hnrss, and lobsters sources
type FeedSource struct {
	name           string
	url            string
	sourceType     string
	useGUID        bool // Use GUID instead of Link (for HN, Lobsters)
	ignoreDays     bool // Set IgnoreDays on items (for Reddit timeless sources)
	isChronological bool // Set IsChronological on items (for "new" sorted feeds)
	parser         *gofeed.Parser
}

// NewFeedSource creates a new feed source
func NewFeedSource(name, url, sourceType string, useGUID, ignoreDays, isChronological bool) *FeedSource {
	return &FeedSource{
		name:           name,
		url:            url,
		sourceType:     sourceType,
		useGUID:        useGUID,
		ignoreDays:     ignoreDays,
		isChronological: isChronological,
		parser:         gofeed.NewParser(),
	}
}

// Fetch retrieves items from the feed
func (f *FeedSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := f.parser.ParseURLWithContext(f.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %w", f.name, err)
	}

	items := make([]models.Item, 0, len(feed.Items))
	for _, item := range feed.Items {
		published := item.PublishedParsed
		if published == nil {
			published = item.UpdatedParsed
		}
		if published == nil {
			continue // Skip items without dates
		}

		var author string
		if item.Author != nil {
			author = item.Author.Name
		}

		content := item.Content
		if content == "" {
			content = item.Description
		}

		link := item.Link
		if f.useGUID && item.GUID != "" {
			link = item.GUID
		}

		items = append(items, models.Item{
			Title:           item.Title,
			Link:            link,
			Description:     item.Description,
			Content:         content,
			Author:          author,
			Published:       *published,
			SourceName:      f.name,
			SourceType:      f.sourceType,
			IgnoreDays:      f.ignoreDays,
			IsChronological: f.isChronological,
		})
	}

	return items, nil
}

// Name returns the source name
func (f *FeedSource) Name() string {
	return f.name
}

// Type returns the source type
func (f *FeedSource) Type() string {
	return f.sourceType
}
