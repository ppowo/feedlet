package source

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ppowo/feedlet/internal/models"
)

// FourPlebsSource implements the Source interface for 4plebs archive
type FourPlebsSource struct {
	name  string
	board string
	limit int
	nsfw  bool
}

// FourPlebsPost represents a post from the 4plebs API
type FourPlebsPost struct {
	DocID            string      `json:"doc_id"`
	Num              string      `json:"num"`
	Subnum           string      `json:"subnum"`
	ThreadNum        string      `json:"thread_num"`
	OP               string      `json:"op"`
	Timestamp        int64       `json:"timestamp"`
	TimestampExpired any `json:"timestamp_expired"` // Can be string or number
	Capcode          string      `json:"capcode"`
	Email            string      `json:"email"`
	Name             string      `json:"name"`
	Trip             string      `json:"trip"`
	Title            string      `json:"title"`
	Comment          string      `json:"comment"`
	PosterCountry    string      `json:"poster_country"`
	Sticky           string      `json:"sticky"`
	Locked           string      `json:"locked"`
	Deleted          any `json:"deleted"` // Can be string or number
	NReplies         *int        `json:"nreplies"`
	NImages          *int        `json:"nimages"`
	FourchanDate     string      `json:"fourchan_date"`
	CommentSanitized string      `json:"comment_sanitized"`
}

// NewFourPlebsSource creates a new 4plebs source
func NewFourPlebsSource(name, board string, limit int, nsfw bool) *FourPlebsSource {
	return &FourPlebsSource{
		name:  name,
		board: board,
		limit: limit,
		nsfw:  nsfw,
	}
}

// Fetch retrieves /film/ threads from 4plebs archive
func (f *FourPlebsSource) Fetch(ctx context.Context) ([]models.Item, error) {
	// Use the search API with type=op to get only OP posts (page 1 only to respect API limits)
	searchURL := fmt.Sprintf("https://archive.4plebs.org/_/api/chan/search/?subject=/film/&boards=%s&type=op&page=1", f.board)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent as required by the API documentation
	req.Header.Set("User-Agent", "feedlet/1.0 (contact: admin@example.com)")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch 4plebs: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var rawResponse map[string]json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&rawResponse); err != nil {
		return nil, fmt.Errorf("failed to decode JSON response: %w", err)
	}

	items := make([]models.Item, 0)

	// Iterate through all numeric keys (0, 1, 2, etc.)
	for key, value := range rawResponse {
		if key == "meta" {
			continue
		}

		var resultSet struct {
			Posts []FourPlebsPost `json:"posts"`
		}
		if err := json.Unmarshal(value, &resultSet); err != nil {
			continue // Skip invalid entries
		}

		for _, post := range resultSet.Posts {
			// Only include threads with '/film/' exactly in the title
			if !strings.Contains(post.Title, "/film/") {
				continue
			}

			// Skip if the thread is deleted (handle both string and number)
			deletedStr := fmt.Sprintf("%v", post.Deleted)
			if deletedStr == "1" {
				continue
			}

			// Only include threads with 150+ replies
			replyCount := 0
			if post.NReplies != nil {
				replyCount = *post.NReplies
			}
			if replyCount < 150 {
				continue
			}

			// Extract title from first two lines of the opening post
			title := extractTitleFromComment(post.CommentSanitized)
			if title == "" {
				title = "/film/ - Thread"
			}

			// Build thread URL
			threadURL := fmt.Sprintf("https://archive.4plebs.org/%s/thread/%s/#%s", f.board, post.ThreadNum, post.ThreadNum)

			// Parse timestamp
			published := time.Unix(post.Timestamp, 0)

			// Clean up the comment for description
			description := strings.TrimSpace(post.CommentSanitized)
			if len(description) > 500 {
				description = description[:500] + "..."
			}

			items = append(items, models.Item{
				Title:       fmt.Sprintf("%s [%d replies]", title, replyCount),
				Link:        threadURL,
				Description: description,
				Content:     post.CommentSanitized,
				Author:      post.Name,
				Published:   published,
				SourceName:  f.name,
				SourceType:  "4plebs",
				IgnoreDays:  true, // Ignore day filtering for /film/ threads
				NSFW:        f.nsfw,
			})

			// Check limit after successful addition
			if f.limit > 0 && len(items) >= f.limit {
				return items, nil
			}
		}
	}

	return items, nil
}

// extractTitleFromComment extracts the title from the first two non-empty lines of a comment
func extractTitleFromComment(comment string) string {
	lines := strings.Split(comment, "\n")

	var firstLines []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			firstLines = append(firstLines, line)
			if len(firstLines) >= 2 {
				break
			}
		}
	}

	if len(firstLines) == 0 {
		return ""
	}

	if len(firstLines) == 1 {
		return firstLines[0]
	}

	// Join the first two lines
	return strings.TrimSpace(firstLines[0] + " " + firstLines[1])
}

// Name returns the source name
func (f *FourPlebsSource) Name() string {
	return f.name
}

// Type returns the source type
func (f *FourPlebsSource) Type() string {
	return "4plebs"
}
