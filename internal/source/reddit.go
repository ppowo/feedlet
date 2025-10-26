package source

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
)

// RedditSource implements the Source interface for Reddit RSS feeds
type RedditSource struct {
	name       string
	url        string
	parser     *gofeed.Parser
	ignoreDays bool
}

// NewRedditSource creates a new Reddit source
func NewRedditSource(name, url string, ignoreDays bool) *RedditSource {
	return &RedditSource{
		name:       name,
		url:        url,
		parser:     gofeed.NewParser(),
		ignoreDays: ignoreDays,
	}
}

// Fetch retrieves items from the Reddit RSS feed
func (r *RedditSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := r.parser.ParseURLWithContext(r.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Reddit feed %s: %w", r.name, err)
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

		// Reddit RSS uses Link for the post URL
		link := item.Link

		items = append(items, models.Item{
			Title:       item.Title,
			Link:        link,
			Description: item.Description,
			Content:     content,
			Author:      author,
			Published:   *published,
			SourceName:  r.name,
			SourceType:  "reddit",
			IgnoreDays:  r.ignoreDays,
		})
	}

	return items, nil
}

// Name returns the source name
func (r *RedditSource) Name() string {
	return r.name
}

// Type returns the source type
func (r *RedditSource) Type() string {
	return "reddit"
}
