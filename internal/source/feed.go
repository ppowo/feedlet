package source

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/mmcdole/gofeed"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

// FeedSource implements the Source interface for RSS/Atom feeds
// Consolidates rss, reddit, hnrss, and lobsters sources
type FeedSource struct {
	name            string
	url             string
	sourceType      string
	useGUID         bool // Use GUID instead of Link (for HN, Lobsters)
	ignoreDays      bool // Set IgnoreDays on items (for Reddit timeless sources)
	isChronological bool // Set IsChronological on items (for "new" sorted feeds)
	parser          *gofeed.Parser
}

// NewFeedSource creates a new feed source
func NewFeedSource(name, url, sourceType string, useGUID, ignoreDays, isChronological bool) *FeedSource {
	return &FeedSource{
		name:            name,
		url:             url,
		sourceType:      sourceType,
		useGUID:         useGUID,
		ignoreDays:      ignoreDays,
		isChronological: isChronological,
		parser:          gofeed.NewParser(),
	}
}

// fetchFeed fetches the feed URL with a proper User-Agent and returns the parsed feed.
func (f *FeedSource) fetchFeed(ctx context.Context) (*gofeed.Feed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, f.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request for %s: %w", f.name, err)
	}
	req.Header.Set("User-Agent", httpclient.RandomUserAgent())

	client := httpclient.GetClient()
	resp, err := client.StandardClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed %s: %w", f.name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch feed %s: http %d %s", f.name, resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read feed %s: %w", f.name, err)
	}

	feed, err := f.parser.ParseString(string(body))
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed %s: %w", f.name, err)
	}
	return feed, nil
}

// Fetch retrieves items from the feed
func (f *FeedSource) Fetch(ctx context.Context) ([]models.Item, error) {
	feed, err := f.fetchFeed(ctx)
	if err != nil {
		return nil, err
	}

	items := make([]models.Item, 0, len(feed.Items))
	for _, item := range feed.Items {
		published := item.PublishedParsed
		if published == nil {
			published = item.UpdatedParsed
		}
		if published == nil {
			continue
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
