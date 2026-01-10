package source

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

// ChanArchiveSource implements the Source interface for 4chan archive APIs
// Supports 4plebs and desuarchive
type ChanArchiveSource struct {
	name        string
	board       string
	limit       int
	nsfw        bool
	archiveType string // "4plebs" or "desuarchive"
	subject     string // Search subject (e.g., "/film/", "/ptg/")
	baseURL     string
	minReplies  int           // Minimum reply count (0 = no filter)
	minAge      time.Duration // Minimum age (0 = no filter)
}

// ChanArchivePost represents a post from the archive API
type ChanArchivePost struct {
	DocID            string `json:"doc_id"`
	Num              string `json:"num"`
	Subnum           string `json:"subnum"`
	ThreadNum        string `json:"thread_num"`
	OP               string `json:"op"`
	Timestamp        int64  `json:"timestamp"`
	TimestampExpired any    `json:"timestamp_expired"`
	Capcode          string `json:"capcode"`
	Email            string `json:"email"`
	Name             string `json:"name"`
	Trip             string `json:"trip"`
	Title            string `json:"title"`
	Comment          string `json:"comment"`
	PosterCountry    string `json:"poster_country"`
	Sticky           string `json:"sticky"`
	Locked           string `json:"locked"`
	Deleted          any    `json:"deleted"`
	NReplies         *int   `json:"nreplies"`
	NImages          *int   `json:"nimages"`
	FourchanDate     string `json:"fourchan_date"`
	CommentSanitized string `json:"comment_sanitized"`
}

// NewFourPlebsSource creates a 4plebs source for /film/
func NewFourPlebsSource(name, board string, limit int, nsfw bool) *ChanArchiveSource {
	return &ChanArchiveSource{
		name:        name,
		board:       board,
		limit:       limit,
		nsfw:        nsfw,
		archiveType: "4plebs",
		subject:     "/film/",
		baseURL:     "https://archive.4plebs.org",
		minReplies:  150,
		minAge:      0,
	}
}

// NewDesuArchiveSource creates a desuarchive source for /ptg/
func NewDesuArchiveSource(name, board string, limit int, nsfw bool) *ChanArchiveSource {
	return &ChanArchiveSource{
		name:        name,
		board:       board,
		limit:       limit,
		nsfw:        nsfw,
		archiveType: "desuarchive",
		subject:     "/ptg/",
		baseURL:     "https://desuarchive.org",
		minReplies:  0,
		minAge:      24 * time.Hour,
	}
}

func (c *ChanArchiveSource) Fetch(ctx context.Context) ([]models.Item, error) {
	searchURL := fmt.Sprintf("%s/_/api/chan/search/?subject=%s&boards=%s&type=op&page=1",
		c.baseURL, c.subject, c.board)

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "feedlet/1.0 (contact: admin@example.com)")
	req.Header.Set("Accept", "application/json")

	resp, err := httpclient.GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", c.archiveType, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawResponse map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	items := make([]models.Item, 0)

	for key, value := range rawResponse {
		if key == "meta" {
			continue
		}

		var resultSet struct {
			Posts []ChanArchivePost `json:"posts"`
		}
		if err := json.Unmarshal(value, &resultSet); err != nil {
			continue
		}

		for _, post := range resultSet.Posts {
			// Filter by subject in title
			if !strings.Contains(post.Title, c.subject) {
				continue
			}

			// Skip deleted posts
			if fmt.Sprintf("%v", post.Deleted) == "1" {
				continue
			}

			published := time.Unix(post.Timestamp, 0)

			// Apply minimum age filter
			if c.minAge > 0 && time.Since(published) <= c.minAge {
				continue
			}

			// Apply minimum replies filter
			replyCount := 0
			if post.NReplies != nil {
				replyCount = *post.NReplies
			}
			if c.minReplies > 0 && replyCount < c.minReplies {
				continue
			}

			// Extract title from comment
			title := extractFirstLines(post.CommentSanitized, 2)
			if title == "" {
				title = c.subject + " - Thread"
			}

			// Add reply count suffix for 4plebs
			if c.archiveType == "4plebs" && replyCount > 0 {
				title = fmt.Sprintf("%s [%d replies]", title, replyCount)
			}

			threadURL := fmt.Sprintf("%s/%s/thread/%s/#%s",
				c.baseURL, c.board, post.ThreadNum, post.ThreadNum)

			description := strings.TrimSpace(post.CommentSanitized)
			if len(description) > 500 {
				description = description[:500] + "..."
			}

			items = append(items, models.Item{
				Title:       title,
				Link:        threadURL,
				Description: description,
				Content:     post.CommentSanitized,
				Author:      post.Name,
				Published:   published,
				SourceName:  c.name,
				SourceType:  c.archiveType,
				IgnoreDays:  true,
				NSFW:        c.nsfw,
			})

			if c.limit > 0 && len(items) >= c.limit {
				return items, nil
			}
		}
	}

	return items, nil
}

// extractFirstLines extracts the first n non-empty lines from a comment
func extractFirstLines(comment string, n int) string {
	lines := strings.Split(comment, "\n")
	var result []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			result = append(result, line)
			if len(result) >= n {
				break
			}
		}
	}

	return strings.Join(result, " ")
}

func (c *ChanArchiveSource) Name() string {
	return c.name
}

func (c *ChanArchiveSource) Type() string {
	return c.archiveType
}
