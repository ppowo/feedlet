package models

import "time"

// Item represents a single feed item from any source
type Item struct {
	Title       string
	Link        string
	Description string
	Content     string
	Author      string
	Published   time.Time
	SourceName  string
	SourceType  string // "rss", "reddit", etc.
	IgnoreDays  bool   // If true, don't filter by age
	NSFW        bool   // If true, apply NSFW styling (red tint)
	HideDate    bool   // If true, don't display the date on frontend
}

// Feed represents a collection of items from all sources
type Feed struct {
	Items     []Item
	UpdatedAt time.Time
	Errors    map[string]string // source name -> error message
}

// SourceConfig represents configuration for a single source
type SourceConfig struct {
	Name           string `yaml:"name"`
	Type           string `yaml:"type"` // "rss", "reddit", etc.
	URL            string `yaml:"url"`
	Interval       int    `yaml:"interval"`        // Base interval in seconds
	IntervalJitter int    `yaml:"interval_jitter"` // Random jitter in seconds (0 to this value)
	IgnoreDays     bool   `yaml:"ignore_days"`     // Don't filter by age
	NSFW           bool   `yaml:"nsfw"`            // Apply NSFW styling (red tint)
	Limit          int    `yaml:"limit"`           // Max items to show (for sources that support it)
	Days           int    `yaml:"days"`            // Number of days to filter (default: 2)
}

// Config represents the application configuration
type Config struct {
	Sources         []SourceConfig `yaml:"sources"`
	Port            int            `yaml:"port"`
	MinFetchInterval int           `yaml:"min_fetch_interval"` // Minimum interval between fetches in seconds (0 = no limit)
	MaxSubscribers  int            `yaml:"max_subscribers"`    // Maximum concurrent subscribers (0 = unlimited)
}
