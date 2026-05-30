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
	HomeURL         string `yaml:"home_url"`
	Interval        int    `yaml:"interval"`
	IntervalJitter  int    `yaml:"interval_jitter"`
	NSFW            bool   `yaml:"nsfw"`
}
type Config struct {
	Port               int             `yaml:"port"`
	MinFetchInterval   int             `yaml:"min_fetch_interval"`
	MaxSubscribers     int             `yaml:"max_subscribers"`
	Sources            []SourceConfig  `yaml:"sources"`
}
