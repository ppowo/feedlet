package config

import (
	"github.com/ppowo/feedlet/internal/models"
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
				Interval:       1800,
				IntervalJitter: 120,
				Days:           6,
			},
			{
				Name:           "r/Italia Personal Finance",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/ItaliaPersonalFinance/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           8,
			},
			{
				Name:           "r/Programming",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/programming/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           8,
			},
			{
				Name:           "HN 350+",
				Type:           "hnrss",
				URL:            "https://hnrss.org/newest?points=350",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           1,
			},
			{
				Name:           "r/Trackers",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/trackers/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           8,
			},
			{
				Name:           "Lobsters Top (1w)",
				Type:           "lobsters",
				URL:            "https://lobste.rs/top/1w.rss",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           2,
			},
			{
				Name:           "r/40k Lore",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/40kLore/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           9,
			},
			{
				Name:           "4plebs /tv/ /film/",
				Type:           "4plebs",
				URL:            "tv",
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
				Interval:       1800,
				IntervalJitter: 120,
				IgnoreDays:     true,
				NSFW:           true,
				Limit:          4,
			},
			{
				Name:           "Meltzer 5.5★+ Matches",
				Type:           "wikipedia",
				URL:            "https://en.wikipedia.org/wiki/List_of_professional_wrestling_matches_rated_5_or_more_stars_by_Dave_Meltzer",
				Interval:       3600,
				IntervalJitter: 300,
				IgnoreDays:     true,
				Limit:          4,
			},
			{
				Name:           "r/Grimdank",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/Grimdank/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           7,
			},
			{
				Name:           "r/Selfhosted",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/selfhosted/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           7,
			},
		},
	}
}
