package models

import "time"

// Item represents a single feed item from any source
type Item struct {
	Title           string
	Link            string
	Description     string
	Content         string
	Author          string
	Published       time.Time
	SourceName      string
	SourceType      string
	IgnoreDays      bool
	NSFW            bool
	IsChronological bool
	HideDate        bool
}

// SourceState represents the runtime health of a source.
type SourceState struct {
	Name                string
	Type                string
	Host                string
	LastAttemptAt       time.Time
	LastSuccessAt       time.Time
	LastError           string
	ConsecutiveFailures int
	Stale               bool
}

// Feed represents a collection of items from all sources.
type Feed struct {
	Items        []Item
	UpdatedAt    time.Time
	Errors       map[string]string
	SourceStates map[string]SourceState
}

// SourceConfig represents configuration for a single source.
type SourceConfig struct {
	Name            string `yaml:"name"`
	Type            string `yaml:"type"`
	URL             string `yaml:"url"`
	Interval        int    `yaml:"interval"`
	IntervalJitter  int    `yaml:"interval_jitter"`
	IgnoreDays      bool   `yaml:"ignore_days"`
	NSFW            bool   `yaml:"nsfw"`
	IsChronological bool   `yaml:"is_chronological"`
	Limit           int    `yaml:"limit"`
	Days            int    `yaml:"days"`
}

// Config represents the application configuration.
type Config struct {
	Sources            []SourceConfig `yaml:"sources"`
	Port               int            `yaml:"port"`
	MinFetchInterval   int            `yaml:"min_fetch_interval"`
	MaxSubscribers     int            `yaml:"max_subscribers"`
	DefaultSourceLimit int            `yaml:"default_source_limit"`
}
