package source

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
)

// RSSSource implements the Source interface for RSS/Atom feeds
type RSSSource struct {
	name   string
	url    string
	parser *gofeed.Parser
}

// NewRSSSource creates a new RSS source
func NewRSSSource(name, url string) *RSSSource {
	return &RSSSource{
		name:   name,
		url:    url,
		parser: gofeed.NewParser(),
	}
}

// Fetch retrieves items from the RSS feed
func (r *RSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := r.parser.ParseURLWithContext(r.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse RSS feed %s: %w", r.name, err)
	}

	items := make([]models.Item, 0, len(feed.Items))
	for _, item := range feed.Items {
		var published = item.PublishedParsed
		if published == nil {
			published = item.UpdatedParsed
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

		items = append(items, models.Item{
			Title:       item.Title,
			Link:        link,
			Description: item.Description,
			Content:     content,
			Author:      author,
			Published:   *published,
			SourceName:  r.name,
			SourceType:  "rss",
		})
	}

	return items, nil
}

// Name returns the source name
func (r *RSSSource) Name() string {
	return r.name
}

// Type returns the source type
func (r *RSSSource) Type() string {
	return "rss"
}
