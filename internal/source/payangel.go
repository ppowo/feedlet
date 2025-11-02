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
	name       string
	url        string
	sections   []string // List of sections to fetch: "FleshSim", "SHS", "SHS-Overtime"
	limit      int      // Total limit across all sections
	ignoreDays bool
}

// NewPayAngelSource creates a new PayAngel source
func NewPayAngelSource(name, url string, sections []string, limit int) *PayAngelSource {
	return &PayAngelSource{
		name:       name,
		url:        url,
		sections:   sections,
		limit:      limit,
		ignoreDays: true, // Archive entries should ignore day filtering
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

	allItems := make([]models.Item, 0)

	// Calculate limit per section (divide total limit by number of sections)
	limitPerSection := p.limit / len(p.sections)
	if limitPerSection == 0 {
		limitPerSection = 1 // At least 1 item per section
	}

	// Fetch from each section
	for _, section := range p.sections {
		var sectionID string
		var sectionPrefix string

		switch section {
		case "FleshSim":
			sectionID = "#fleshsim-contact-for-information"
			sectionPrefix = "FleshSim"
		case "SHS":
			sectionID = "#shs-the-sam-hyde-show"
			sectionPrefix = "SHS"
		case "SHS-Overtime":
			sectionID = "#shs-overtime"
			sectionPrefix = "SHS - Overtime"
		default:
			continue // Skip unknown sections
		}

		items := p.fetchSection(doc, sectionID, sectionPrefix, limitPerSection)
		allItems = append(allItems, items...)
	}

	return allItems, nil
}

// fetchSection extracts items from a specific section of the page
func (p *PayAngelSource) fetchSection(doc *goquery.Document, sectionID, sectionPrefix string, limit int) []models.Item {
	items := make([]models.Item, 0)

	// Find the section heading
	section := doc.Find(sectionID).First()
	if section.Length() == 0 {
		return items // Return empty if section not found
	}

	// Regex to extract episode information from text like "EP01" or "Episode 1"
	epRegex := regexp.MustCompile(`(?i)(EP|Episode)\s*(\d+)`)

	// Navigate through siblings to find episode entries
	// Episodes are in <p> tags followed by code blocks with magnet links
	// Start from the section heading itself
	current := section

	for current.Next().Length() > 0 {
		current = current.Next()

		// Stop when we hit the next heading
		if current.Is("h1, h2, h3, h4") {
			break
		}

		// Look for paragraph tags that contain episode titles
		if current.Is("p") {
			text := strings.TrimSpace(current.Text())

			// Skip empty paragraphs or paragraphs that don't look like episode titles
			if text == "" || len(text) > 200 {
				continue
			}

			// Extract episode number if present
			epMatch := epRegex.FindStringSubmatch(text)
			var episodeNum string
			if len(epMatch) >= 3 {
				episodeNum = epMatch[2]
			}

			// Clean up the title (remove trailing <br>, etc.)
			title := strings.TrimSuffix(text, "\n")
			title = strings.TrimSpace(title)

			// Skip if this doesn't look like an episode entry
			if title == "" {
				continue
			}

			// Format title as "SECTION - Episode Title"
			formattedTitle := fmt.Sprintf("%s - %s", sectionPrefix, title)

			// Create description with episode number if available
			description := "Sam Hyde Archive"
			if episodeNum != "" {
				description = fmt.Sprintf("Episode %s", episodeNum)
			}

			// Set a fake old date (10 years ago) to ensure this appears last on the frontend
			oldDate := time.Now().AddDate(-10, 0, 0)

			item := models.Item{
				Title:       formattedTitle,
				Link:        p.url, // Link to the archive page, not the magnet link
				Description: description,
				Content:     fmt.Sprintf("Episode from %s\n\nView at: %s", sectionPrefix, p.url),
				Author:      "Sam Hyde Archive",
				Published:   oldDate, // Use old date to appear last
				SourceName:  p.name,
				SourceType:  "payangel",
				IgnoreDays:  p.ignoreDays,
				HideDate:    true, // Don't show date on frontend
			}

			items = append(items, item)

			// Stop if we've reached the limit for this section
			if limit > 0 && len(items) >= limit {
				break
			}
		}
	}

	return items
}

// Name returns the source name
func (p *PayAngelSource) Name() string {
	return p.name
}

// Type returns the source type
func (p *PayAngelSource) Type() string {
	return "payangel"
}
