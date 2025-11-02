package source

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/ppowo/feedlet/internal/models"
)

// PayAngelSource implements the Source interface for the Sam Hyde archive
type PayAngelSource struct {
	name     string
	url      string
	sections []string // List of sections: "FleshSim", "SHS", "SHS-Overtime"
	limit    int      // Total limit across all sections
}

// NewPayAngelSource creates a new PayAngel source
func NewPayAngelSource(name, url string, sections []string, limit int) *PayAngelSource {
	return &PayAngelSource{
		name:     name,
		url:      url,
		sections: sections,
		limit:    limit,
	}
}

// Fetch retrieves episode data from the PayAngel archive
func (p *PayAngelSource) Fetch(ctx context.Context) ([]models.Item, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Feedlet/1.0 (RSS aggregator)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payangel: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Map section identifiers to their HTML IDs and display names
	sectionMap := map[string]struct{ id, prefix string }{
		"FleshSim":     {"#fleshsim-contact-for-information", "FleshSim"},
		"SHS":          {"#shs-the-sam-hyde-show", "SHS"},
		"SHS-Overtime": {"#shs-overtime", "SHS - Overtime"},
	}

	limitPerSection := p.limit / len(p.sections)
	if limitPerSection == 0 {
		limitPerSection = 1
	}

	// Fetch latest N items from each section
	sectionItems := make([][]models.Item, 0, len(p.sections))
	for _, section := range p.sections {
		info, exists := sectionMap[section]
		if !exists {
			continue
		}
		items := p.fetchLatestFromSection(doc, info.id, info.prefix, limitPerSection)
		if len(items) > 0 {
			sectionItems = append(sectionItems, items)
		}
	}

	// Interleave items: latest from each section first, then second-latest, etc.
	return p.interleaveAndTimestamp(sectionItems), nil
}

// fetchLatestFromSection extracts the latest N episodes from a section
func (p *PayAngelSource) fetchLatestFromSection(doc *goquery.Document, sectionID, sectionPrefix string, limit int) []models.Item {
	section := doc.Find(sectionID).First()
	if section.Length() == 0 {
		return nil
	}

	epRegex := regexp.MustCompile(`(?i)(EP|Episode)\s*(\d+)`)
	allEpisodes := make([]models.Item, 0)

	// Traverse sibling elements looking for episode paragraphs
	for el := section.Next(); el.Length() > 0; el = el.Next() {
		if el.Is("h1, h2, h3, h4") {
			break // Hit next section
		}

		if !el.Is("p") {
			continue
		}

		text := strings.TrimSpace(el.Text())
		if text == "" || len(text) > 200 {
			continue // Skip empty or overly long paragraphs
		}

		// Extract episode number for description
		episodeNum := ""
		if match := epRegex.FindStringSubmatch(text); len(match) >= 3 {
			episodeNum = match[2]
		}

		description := "Sam Hyde Archive"
		if episodeNum != "" {
			description = fmt.Sprintf("Episode %s", episodeNum)
		}

		allEpisodes = append(allEpisodes, models.Item{
			Title:       fmt.Sprintf("%s - %s", sectionPrefix, text),
			Link:        p.url,
			Description: description,
			Author:      "Sam Hyde Archive",
			SourceName:  p.name,
			SourceType:  "payangel",
			IgnoreDays:  true,
			HideDate:    true,
		})
	}

	// Return the last N episodes (latest)
	if len(allEpisodes) <= limit {
		return allEpisodes
	}
	return allEpisodes[len(allEpisodes)-limit:]
}

// interleaveAndTimestamp interleaves items from multiple sections and assigns timestamps
func (p *PayAngelSource) interleaveAndTimestamp(sectionItems [][]models.Item) []models.Item {
	if len(sectionItems) == 0 {
		return nil
	}

	// Find max items per section
	maxItems := 0
	for _, items := range sectionItems {
		if len(items) > maxItems {
			maxItems = len(items)
		}
	}

	// Interleave: iterate from last to first (latest episodes first)
	result := make([]models.Item, 0, p.limit)
	baseDate := time.Now().AddDate(-10, 0, 0) // 10 years ago to appear last on dashboard

	for i := maxItems - 1; i >= 0; i-- {
		for _, items := range sectionItems {
			if i < len(items) {
				item := items[i]
				// Assign staggered timestamps to preserve order
				item.Published = baseDate.Add(time.Duration(-len(result)) * time.Hour)
				result = append(result, item)
			}
		}
	}

	return result
}

// Name returns the source name
func (p *PayAngelSource) Name() string {
	return p.name
}

// Type returns the source type
func (p *PayAngelSource) Type() string {
	return "payangel"
}
