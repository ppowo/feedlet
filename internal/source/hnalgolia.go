package source

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

const (
	hnAlgoliaBaseURL     = "https://hn.algolia.com/api/v1/search_by_date"
	defaultHNItemsPerReq = 20
	maxHNItemsPerReq     = 100
)

var hnFeedTagByPath = map[string]string{
	"ask":         "ask_hn",
	"frontpage":   "front_page",
	"newcomments": "comment",
	"newest":      "story",
	"polls":       "poll",
	"show":        "show_hn",
}

// HNAlgoliaSource fetches Hacker News items directly from the Algolia API.
//
// It accepts direct Algolia URLs and simplified HN feed-style URLs for the
// common feed types used by Feedlet.
type HNAlgoliaSource struct {
	name       string
	url        string
	sourceType string
}

// NewHNAlgoliaSource creates a new direct Hacker News Algolia source.
func NewHNAlgoliaSource(name, rawURL, sourceType string) *HNAlgoliaSource {
	return &HNAlgoliaSource{
		name:       name,
		url:        rawURL,
		sourceType: sourceType,
	}
}

type hnAlgoliaResponse struct {
	Hits []hnAlgoliaHit `json:"hits"`
}

type hnAlgoliaHit struct {
	Author      string `json:"author"`
	CommentText string `json:"comment_text"`
	CreatedAt   string `json:"created_at"`
	CreatedAtI  int64  `json:"created_at_i"`
	NumComments int    `json:"num_comments"`
	ObjectID    string `json:"objectID"`
	Points      int    `json:"points"`
	StoryText   string `json:"story_text"`
	StoryTitle  string `json:"story_title"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

func (h *HNAlgoliaSource) Fetch(ctx context.Context) ([]models.Item, error) {
	requestURL, includeDescription, err := h.buildRequestURL()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create HN Algolia request for %s: %w", h.name, err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", httpclient.RandomUserAgent())

	client := httpclient.GetClient()
	resp, err := client.StandardClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch HN Algolia data for %s: %w", h.name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch HN Algolia data for %s: http %d %s", h.name, resp.StatusCode, resp.Status)
	}

	var payload hnAlgoliaResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, fmt.Errorf("failed to decode HN Algolia response for %s: %w", h.name, err)
	}

	items := make([]models.Item, 0, len(payload.Hits))
	for _, hit := range payload.Hits {
		published, err := hit.publishedAt()
		if err != nil {
			continue
		}

		title := strings.TrimSpace(hit.Title)
		if title == "" {
			title = strings.TrimSpace(hit.StoryTitle)
		}
		if title == "" {
			title = h.commentLink(hit.ObjectID)
		}

		description := ""
		if includeDescription {
			description = h.buildDescription(hit)
		}

		content := strings.TrimSpace(hit.StoryText)
		if content == "" {
			content = strings.TrimSpace(hit.CommentText)
		}
		if content == "" {
			content = description
		}

		items = append(items, models.Item{
			Title:           title,
			Link:            h.commentLink(hit.ObjectID),
			Description:     description,
			Content:         content,
			Author:          hit.Author,
			Published:       published,
			SourceName:      h.name,
			SourceType:      h.sourceType,
			IgnoreDays:      false,
			IsChronological: false,
		})
	}

	return items, nil
}

func (h *HNAlgoliaSource) Name() string {
	return h.name
}

func (h *HNAlgoliaSource) Type() string {
	return h.sourceType
}

func (h *HNAlgoliaSource) buildRequestURL() (string, bool, error) {
	srcURL, err := url.Parse(h.url)
	if err != nil {
		return "", false, fmt.Errorf("invalid HN source URL for %s: %w", h.name, err)
	}

	if strings.EqualFold(srcURL.Host, "hn.algolia.com") {
		q := srcURL.Query()
		if q.Get("hitsPerPage") == "" {
			q.Set("hitsPerPage", strconv.Itoa(defaultHNItemsPerReq))
		}
		srcURL.RawQuery = q.Encode()
		return srcURL.String(), shouldIncludeDescription(srcURL.Query()), nil
	}

	path := strings.Trim(srcURL.Path, "/")
	tags, ok := hnFeedTagByPath[path]
	if !ok {
		return "", false, fmt.Errorf("unsupported HN feed path %q for %s", path, h.name)
	}

	sourceQuery := srcURL.Query()
	requestQuery := url.Values{}
	requestQuery.Set("tags", tags)
	requestQuery.Set("hitsPerPage", strconv.Itoa(parseBoundedInt(sourceQuery.Get("count"), defaultHNItemsPerReq, 1, maxHNItemsPerReq)))

	if q := strings.TrimSpace(sourceQuery.Get("q")); q != "" {
		requestQuery.Set("query", q)
	}

	numericFilters := make([]string, 0, 2)
	if points, ok := parsePositiveInt(sourceQuery.Get("points")); ok {
		numericFilters = append(numericFilters, fmt.Sprintf("points>%d", points))
	}
	if comments, ok := parsePositiveInt(sourceQuery.Get("comments")); ok {
		numericFilters = append(numericFilters, fmt.Sprintf("num_comments>%d", comments))
	}
	if len(numericFilters) > 0 {
		requestQuery.Set("numericFilters", strings.Join(numericFilters, ","))
	}

	return hnAlgoliaBaseURL + "?" + requestQuery.Encode(), shouldIncludeDescription(sourceQuery), nil
}

func (h *HNAlgoliaSource) commentLink(objectID string) string {
	return "https://news.ycombinator.com/item?id=" + objectID
}

func (h *HNAlgoliaSource) buildDescription(hit hnAlgoliaHit) string {
	commentsURL := h.commentLink(hit.ObjectID)
	var b strings.Builder

	storyText := strings.TrimSpace(hit.StoryText)
	if storyText != "" {
		fmt.Fprintf(&b, "\n<p>%s</p>\n", storyText)
	} else if hit.URL != "" {
		escapedURL := html.EscapeString(hit.URL)
		fmt.Fprintf(&b, "\n<p>Article URL: <a href=\"%s\">%s</a></p>\n", escapedURL, escapedURL)
	}

	if b.Len() > 0 {
		b.WriteString("<hr>\n")
	}

	escapedCommentsURL := html.EscapeString(commentsURL)
	fmt.Fprintf(&b, "<p>Comments URL: <a href=\"%s\">%s</a></p>\n", escapedCommentsURL, escapedCommentsURL)
	fmt.Fprintf(&b, "<p>Points: %d</p>\n", hit.Points)
	fmt.Fprintf(&b, "<p># Comments: %d</p>\n", hit.NumComments)

	return b.String()
}

func (h hnAlgoliaHit) publishedAt() (time.Time, error) {
	if h.CreatedAtI > 0 {
		return time.Unix(h.CreatedAtI, 0).UTC(), nil
	}
	if h.CreatedAt == "" {
		return time.Time{}, fmt.Errorf("missing created_at")
	}
	published, err := time.Parse(time.RFC3339, h.CreatedAt)
	if err != nil {
		return time.Time{}, err
	}
	return published, nil
}

func parsePositiveInt(value string) (int, bool) {
	value = strings.TrimSpace(value)
	if value == "" {
		return 0, false
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 0 {
		return 0, false
	}
	return n, true
}

func parseBoundedInt(value string, fallback, min, max int) int {
	if n, ok := parsePositiveInt(value); ok {
		if n < min {
			return min
		}
		if n > max {
			return max
		}
		return n
	}
	return fallback
}

func shouldIncludeDescription(q url.Values) bool {
	return strings.TrimSpace(q.Get("description")) != "0"
}
