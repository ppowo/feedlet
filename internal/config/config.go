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
		DefaultSourceLimit: 4,
		Sources: []models.SourceConfig{
			{
				Name:           "r/Italia Career Advice",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/italiacareeradvice/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/italiacareeradvice/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           10,
			},
			{
				Name:           "r/Italia Personal Finance",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/ItaliaPersonalFinance/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/ItaliaPersonalFinance/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           9,
			},
			{
				Name:           "r/Programming",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/programming/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/programming/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           8,
			},
			{
				Name:           "HN 350+",
				Type:           "hnalgolia",
				URL:            "https://hn.algolia.com/api/v1/search_by_date?tags=story&numericFilters=points%3E350",
				HomeURL:        "https://hckrnews.com/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           3,
			},
			{
				Name:           "r/Trackers",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/trackers/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/trackers/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           10,
			},
			{
			Name:           "Lobsters Top (1d)",
				Type:           "lobsters",
				URL:            "https://lobste.rs/top/1d.rss",
				HomeURL:        "https://lobste.rs/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           4,
			},
			{
				Name:           "r/40k Lore",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/40kLore/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/40kLore/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           9,
			},
			{
				Name:           "4plebs /tv/ /film/",
				Type:           "4plebs",
				URL:            "tv",
				HomeURL:        "",
				Interval:       1800,
				IntervalJitter: 120,
				IgnoreDays:     true,
				NSFW:           true,
				Limit:          4,
			},
			{
				Name:           "DesuArchive /g/ /ptg/",
				Type:           "desuarchive",
				URL:            "g",
				HomeURL:        "",
				Interval:       1800,
				IntervalJitter: 120,
				IgnoreDays:     true,
				NSFW:           true,
				Limit:          4,
			},
			{
				Name:           source.CurrentCagematchSourceName(),
				Type:           "cagematch",
				URL:            "https://www.cagematch.net/",
				HomeURL:        source.CurrentCagematchHomeURL(),
				Interval:       3600,
				IntervalJitter: 300,
				IgnoreDays:     true,
				Limit:          4,
			},
			{
				Name:           "r/kitchencels",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/kitchencels/top/.rss?t=month",
				HomeURL:        "https://old.reddit.com/r/kitchencels/top/",
				Interval:       1800,
				IntervalJitter: 120,
				NSFW:           true,
				Days:           10,
			},
			{
				Name:           "r/SquaredCircle",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/SquaredCircle/top/.rss?t=week",
				HomeURL:        "https://old.reddit.com/r/SquaredCircle/top/",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           3,
			},
		},
	}
}
