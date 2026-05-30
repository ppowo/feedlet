package config

import (
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source"
)

// GetConfig returns the embedded application configuration
func GetConfig() *models.Config {
	return &models.Config{
		Port:               3737,
		MinFetchInterval:   5, // Default 5 second minimum between fetches per source
		MaxSubscribers:     1000,
		Sources: []models.SourceConfig{
			{
				Name:           "r/Italia Career Advice",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/italiacareeradvice/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/italiacareeradvice/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/Italia Personal Finance",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/ItaliaPersonalFinance/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/ItaliaPersonalFinance/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/Programming",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/programming/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/programming/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "HN 350+",
				Type:           "hnalgolia",
				URL:            "https://hn.algolia.com/api/v1/search_by_date?tags=story&numericFilters=points%3E350",
				HomeURL:        "https://hckrnews.com/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/Trackers",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/trackers/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/trackers/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "Tildes ~tech",
				Type:           "tildes",
				URL:            "https://tildes.net/~tech?order=votes&period=90d",
				HomeURL:        "https://tildes.net/~tech?order=votes&period=90d",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/40k Lore",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/40kLore/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/40kLore/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/TrueFilm",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/truefilm/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/truefilm/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "DesuArchive /g/ /ptg/",
				Type:           "desuarchive",
				URL:            "g",
				HomeURL:        "",
				Interval:       1800,
				IntervalJitter: 120,
				NSFW:           true,
			},
			{
				Name:           source.CurrentMeltzerWikiSourceName(),
				Type:           "meltzerwiki",
				URL:            source.CurrentMeltzerWikiHomeURL(),
				HomeURL:        source.CurrentMeltzerWikiHomeURL(),
				Interval:       3600,
				IntervalJitter: 300,
			},
			{
				Name:           "r/technology",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/technology/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/technology/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
			{
				Name:           "r/SquaredCircle",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/SquaredCircle/top/.rss?t=week",
				HomeURL:        "https://old.reddit.com/r/SquaredCircle/top/",
				Interval:       1800,
				IntervalJitter: 120,
			},
		},
	}
}