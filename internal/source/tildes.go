package source

import (
	"context"
	"fmt"
	"net/http"
	neturl "net/url"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

// TildesSource fetches topic listings directly from Tildes HTML pages.
type TildesSource struct {
	name       string
	url        string
	limit      int
	ignoreDays bool
}

// NewTildesSource creates a new Tildes source.
func NewTildesSource(name, rawURL string, limit int, ignoreDays bool) *TildesSource {
	return &TildesSource{
		name:       name,
		url:        normalizeTildesURL(rawURL),
		limit:      limit,
		ignoreDays: ignoreDays,
	}
}

func normalizeTildesURL(rawURL string) string {
	parsed, err := neturl.Parse(rawURL)
	if err != nil {
		return rawURL
	}

	switch {
	case strings.HasSuffix(parsed.Path, "/topics.atom"):
		parsed.Path = strings.TrimSuffix(parsed.Path, "/topics.atom")
	case strings.HasSuffix(parsed.Path, "/topics.rss"):
		parsed.Path = strings.TrimSuffix(parsed.Path, "/topics.rss")
	}

	return parsed.String()
}

func resolveTildesURL(baseURL *neturl.URL, href string) string {
	ref, err := neturl.Parse(strings.TrimSpace(href))
	if err != nil {
		return href
	}
	if baseURL == nil {
		return ref.String()
	}
	return baseURL.ResolveReference(ref).String()
}

func (t *TildesSource) fetchDocument(ctx context.Context, fetchURL string) (*goquery.Document, *neturl.URL, error) {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, fetchURL, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("tildes: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", httpclient.RandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := httpclient.GetClient().Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("tildes: failed to fetch listing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("tildes: unexpected status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("tildes: failed to parse HTML: %w", err)
	}

	return doc, resp.Request.URL, nil
}

func (t *TildesSource) parseTopic(article *goquery.Selection, baseURL *neturl.URL) (models.Item, bool) {
	titleLink := article.Find("h1.topic-title a[href]").First()
	title := strings.TrimSpace(titleLink.Text())
	if title == "" {
		return models.Item{}, false
	}

	linkHref := strings.TrimSpace(article.Find("footer.topic-info .topic-info-comments a[href]").First().AttrOr("href", ""))
	if linkHref == "" {
		linkHref = strings.TrimSpace(titleLink.AttrOr("href", ""))
	}
	if linkHref == "" {
		return models.Item{}, false
	}

	datetime := strings.TrimSpace(article.Find("footer.topic-info time[datetime]").First().AttrOr("datetime", ""))
	if datetime == "" {
		return models.Item{}, false
	}

	published, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return models.Item{}, false
	}

	return models.Item{
		Title:       title,
		Link:        resolveTildesURL(baseURL, linkHref),
		Description: strings.TrimSpace(article.Find("details.topic-text-excerpt summary").Text()),
		Author:      strings.TrimSpace(article.AttrOr("data-topic-posted-by", "")),
		Published:   published,
		SourceName:  t.name,
		SourceType:  "tildes",
		IgnoreDays:  t.ignoreDays,
	}, true
}

func (t *TildesSource) parseListingPage(doc *goquery.Document, baseURL *neturl.URL) ([]models.Item, error) {
	listing := doc.Find("ol.topic-listing")
	if listing.Length() == 0 {
		return nil, fmt.Errorf("tildes: topic listing not found")
	}

	items := make([]models.Item, 0, 32)
	listing.Find("article.topic").Each(func(_ int, article *goquery.Selection) {
		if item, ok := t.parseTopic(article, baseURL); ok {
			items = append(items, item)
		}
	})

	return items, nil
}

// Fetch retrieves topics from the configured Tildes listing.
func (t *TildesSource) Fetch(ctx context.Context) ([]models.Item, error) {
	doc, baseURL, err := t.fetchDocument(ctx, t.url)
	if err != nil {
		return nil, err
	}

	items, err := t.parseListingPage(doc, baseURL)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})

	if t.limit > 0 && len(items) > t.limit {
		items = items[:t.limit]
	}

	return items, nil
}

// Name returns the source name.
func (t *TildesSource) Name() string {
	return t.name
}

// Type returns the source type.
func (t *TildesSource) Type() string {
	return "tildes"
}
