package config

import (
	"github.com/ppowo/feedlet/internal/models"
)

// GetConfig returns the embedded application configuration
func GetConfig() *models.Config {
	return &models.Config{
		Port: 3737,
		Sources: []models.SourceConfig{
			{
				Name:           "r/italiacareeradvice",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/italiacareeradvice/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           6,
			},
			{
				Name:           "r/programming",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/programming/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           7,
			},
			{
				Name:           "Hacker News",
				Type:           "hnrss",
				URL:            "https://hnrss.org/newest?points=350",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           2,
			},
			{
				Name:           "r/thelastpsychiatrist",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/thelastpsychiatrist/top/.rss?sort=top&t=all",
				Interval:       3600,
				IntervalJitter: 300,
				IgnoreDays:     true,
				Limit:          6,
			},
			{
				Name:           "Lobsters",
				Type:           "lobsters",
				URL:            "https://lobste.rs/top/1w.rss",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           2,
			},
			{
				Name:           "r/40kLore",
				Type:           "reddit",
				URL:            "https://old.reddit.com/r/40kLore/top/.rss?t=month",
				Interval:       1800,
				IntervalJitter: 120,
				Days:           9,
			},
			{
				Name:           "4plebs /film/",
				Type:           "4plebs",
				URL:            "tv",
				Interval:       1800,
				IntervalJitter: 120,
				IgnoreDays:     true,
				NSFW:           true,
				Limit:          4,
			},
			{
				Name:           "DesuArchive /ptg/",
				Type:           "desuarchive",
				URL:            "g",
				Interval:       1800,
				IntervalJitter: 120,
				IgnoreDays:     true,
				NSFW:           true,
				Limit:          4,
			},
			{
				Name:           "Meltzer 5★ Matches",
				Type:           "wikipedia",
				URL:            "https://en.wikipedia.org/wiki/List_of_professional_wrestling_matches_rated_5_or_more_stars_by_Dave_Meltzer",
				Interval:       3600,
				IntervalJitter: 300,
				IgnoreDays:     true,
				Limit:          6,
			},
		},
	}
}
