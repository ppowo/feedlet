package source

import (
	"context"
	"fmt"
	"net/http"
	neturl "net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

const (
	tildesDefaultMaxPages = 10
	tildesDefaultMaxItems = 500
)

// TildesSource fetches topic listings directly from Tildes HTML pages.
type TildesSource struct {
	name     string
	url      string
	maxPages int
	maxItems int
}

// NewTildesSource creates a new Tildes source.
func NewTildesSource(name, rawURL string) *TildesSource {
	return &TildesSource{
		name:     name,
		url:      normalizeTildesURL(rawURL),
		maxPages: tildesDefaultMaxPages,
		maxItems: tildesDefaultMaxItems,
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
	title := strings.TrimSpace(article.Find("h1.topic-title a").First().Text())
	if title == "" {
		return models.Item{}, false
	}

	commentsLink, ok := article.Find("footer.topic-info .topic-info-comments a[href]").First().Attr("href")
	if !ok || strings.TrimSpace(commentsLink) == "" {
		return models.Item{}, false
	}

	datetime, ok := article.Find("footer.topic-info time[datetime]").First().Attr("datetime")
	if !ok || strings.TrimSpace(datetime) == "" {
		return models.Item{}, false
	}

	published, err := time.Parse(time.RFC3339, strings.TrimSpace(datetime))
	if err != nil {
		return models.Item{}, false
	}

	return models.Item{
		Title:       title,
		Link:        resolveTildesURL(baseURL, commentsLink),
		Description: strings.TrimSpace(article.Find("details.topic-text-excerpt summary").Text()),
		Author:      strings.TrimSpace(article.AttrOr("data-topic-posted-by", "")),
		Published:   published,
		SourceName:  t.name,
		SourceType:  "tildes",
	}, true
}

func nextTildesPageURL(doc *goquery.Document, baseURL *neturl.URL) string {
	nextURL := ""
	doc.Find(".pagination a.page-item[href]").EachWithBreak(func(_ int, s *goquery.Selection) bool {
		if !strings.EqualFold(strings.TrimSpace(s.Text()), "Next") {
			return true
		}

		href, ok := s.Attr("href")
		if !ok {
			return true
		}

		nextURL = resolveTildesURL(baseURL, href)
		return false
	})
	return nextURL
}

func (t *TildesSource) parseListingPage(doc *goquery.Document, baseURL *neturl.URL) ([]models.Item, string, error) {
	items := make([]models.Item, 0, 32)

	doc.Find("ol.topic-listing article.topic").Each(func(_ int, article *goquery.Selection) {
		if item, ok := t.parseTopic(article, baseURL); ok {
			items = append(items, item)
		}
	})

	return items, nextTildesPageURL(doc, baseURL), nil
}

// Fetch retrieves topics from the configured Tildes listing.
func (t *TildesSource) Fetch(ctx context.Context) ([]models.Item, error) {
	nextURL := t.url
	items := make([]models.Item, 0, 50)
	seen := make(map[string]struct{})
	visited := make(map[string]struct{})

	for page := 0; page < t.maxPages && nextURL != ""; page++ {
		if _, ok := visited[nextURL]; ok {
			break
		}
		visited[nextURL] = struct{}{}

		doc, baseURL, err := t.fetchDocument(ctx, nextURL)
		if err != nil {
			return nil, err
		}

		pageItems, nextPageURL, err := t.parseListingPage(doc, baseURL)
		if err != nil {
			return nil, err
		}

		for _, item := range pageItems {
			if _, ok := seen[item.Link]; ok {
				continue
			}
			seen[item.Link] = struct{}{}
			items = append(items, item)
			if len(items) >= t.maxItems {
				return items, nil
			}
		}

		if len(pageItems) == 0 {
			break
		}
		nextURL = nextPageURL
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
