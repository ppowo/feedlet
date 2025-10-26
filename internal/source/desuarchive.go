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

// DesuArchiveSource implements the Source interface for DesuArchive archive
type DesuArchiveSource struct {
	name  string
	board string
	limit int
	nsfw  bool
}

// DesuArchivePost represents a post from the DesuArchive API
type DesuArchivePost struct {
	DocID            string      `json:"doc_id"`
	Num              string      `json:"num"`
	Subnum           string      `json:"subnum"`
	ThreadNum        string      `json:"thread_num"`
	OP               string      `json:"op"`
	Timestamp        int64       `json:"timestamp"`
	TimestampExpired any         `json:"timestamp_expired"` // Can be string or number
	Capcode          string      `json:"capcode"`
	Email            string      `json:"email"`
	Name             string      `json:"name"`
	Trip             string      `json:"trip"`
	Title            string      `json:"title"`
	Comment          string      `json:"comment"`
	PosterCountry    string      `json:"poster_country"`
	Sticky           string      `json:"sticky"`
	Locked           string      `json:"locked"`
	Deleted          any         `json:"deleted"` // Can be string or number
	NReplies         *int        `json:"nreplies"`
	NImages          *int        `json:"nimages"`
	FourchanDate     string      `json:"fourchan_date"`
	CommentSanitized string      `json:"comment_sanitized"`
}

// NewDesuArchiveSource creates a new DesuArchive source
func NewDesuArchiveSource(name, board string, limit int, nsfw bool) *DesuArchiveSource {
	return &DesuArchiveSource{
		name:  name,
		board: board,
		limit: limit,
		nsfw:  nsfw,
	}
}

// Fetch retrieves /ptg/ threads from DesuArchive archive
func (d *DesuArchiveSource) Fetch(ctx context.Context) ([]models.Item, error) {
	// Use the search API with type=op to get only OP posts (page 1 only to respect API limits)
	searchURL := fmt.Sprintf("https://desuarchive.org/_/api/chan/search/?subject=/ptg/&boards=%s&type=op&page=1", d.board)

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
		return nil, fmt.Errorf("failed to fetch DesuArchive: %w", err)
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
			Posts []DesuArchivePost `json:"posts"`
		}
		if err := json.Unmarshal(value, &resultSet); err != nil {
			continue // Skip invalid entries
		}

		for _, post := range resultSet.Posts {
			// Only include threads with '/ptg/' exactly in the title
			if !strings.Contains(post.Title, "/ptg/") {
				continue
			}

			// Skip if the thread is deleted (handle both string and number)
			deletedStr := fmt.Sprintf("%v", post.Deleted)
			if deletedStr == "1" {
				continue
			}

			// Parse timestamp first to filter by age
			published := time.Unix(post.Timestamp, 0)

			// Only include threads older than 24 hours
			if time.Since(published) <= 24*time.Hour {
				continue
			}

			// Extract title from first line with text from the post
			title := extractFirstLineFromComment(post.CommentSanitized)
			if title == "" {
				title = "/ptg/ - Thread"
			}

			// Build thread URL
			threadURL := fmt.Sprintf("https://desuarchive.org/%s/thread/%s/#%s", d.board, post.ThreadNum, post.ThreadNum)

			// Clean up the comment for description
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
				SourceName:  d.name,
				SourceType:  "desuarchive",
				IgnoreDays:  true, // Ignore day filtering for /ptg/ threads
				NSFW:        d.nsfw,
			})

			// Check limit after successful addition
			if d.limit > 0 && len(items) >= d.limit {
				return items, nil
			}
		}
	}

	return items, nil
}

// extractFirstLineFromComment extracts the first line with text from a comment
func extractFirstLineFromComment(comment string) string {
	lines := strings.Split(comment, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			return line
		}
	}

	return ""
}

// Name returns the source name
func (d *DesuArchiveSource) Name() string {
	return d.name
}

// Type returns the source type
func (d *DesuArchiveSource) Type() string {
	return "desuarchive"
}