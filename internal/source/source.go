package source

import (
	"context"
	"github.com/ppowo/feedlet/internal/models"
)

// Source is the interface that all source types must implement
// This allows easy extension to non-RSS sources (Reddit API, HN API, etc.)
type Source interface {
	// Fetch retrieves items from the source
	Fetch(ctx context.Context) ([]models.Item, error)

	// Name returns the source name
	Name() string

	// Type returns the source type (rss, reddit, etc.)
	Type() string
}
