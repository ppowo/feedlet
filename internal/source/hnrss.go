package source

import (
	"context"
	"fmt"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
)

// HNRSSSource implements the Source interface for hnrss.org feeds
type HNRSSSource struct {
	name   string
	url    string
	parser *gofeed.Parser
}

// NewHNRSSSource creates a new hnrss.org source
func NewHNRSSSource(name, url string) *HNRSSSource {
	return &HNRSSSource{
		name:   name,
		url:    url,
		parser: gofeed.NewParser(),
	}
}

// Fetch retrieves items from the hnrss.org feed
func (h *HNRSSSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := h.parser.ParseURLWithContext(h.url, ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to parse hnrss feed %s: %w", h.name, err)
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

		// For HN, the GUID contains the comments URL
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
			SourceName:  h.name,
			SourceType:  "hnrss",
		})
	}

	return items, nil
}

// Name returns the source name
func (h *HNRSSSource) Name() string {
	return h.name
}

// Type returns the source type
func (h *HNRSSSource) Type() string {
	return "hnrss"
}
