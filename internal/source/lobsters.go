package source

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
)

// LobstersSource implements the Source interface for Lobsters RSS feeds
type LobstersSource struct {
	name   string
	url    string
	parser *gofeed.Parser
}

// NewLobstersSource creates a new Lobsters source
func NewLobstersSource(name, url string) *LobstersSource {
	return &LobstersSource{
		name:   name,
		url:    url,
		parser: gofeed.NewParser(),
	}
}

// Fetch retrieves items from the Lobsters RSS feed
func (l *LobstersSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := l.parser.ParseURLWithContext(l.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Lobsters feed %s: %w", l.name, err)
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

		// For Lobsters, use the GUID which contains the comments link
		link := item.GUID
		if link == "" {
			link = item.Link
		}

		items = append(items, models.Item{
			Title:       item.Title,
			Link:        link,
			Description: item.Description,
			Content:     content,
			Author:      author,
			Published:   *published,
			SourceName:  l.name,
			SourceType:  "lobsters",
		})
	}

	return items, nil
}

// Name returns the source name
func (l *LobstersSource) Name() string {
	return l.name
}

// Type returns the source type
func (l *LobstersSource) Type() string {
	return "lobsters"
}
