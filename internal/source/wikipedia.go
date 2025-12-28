package source

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/araddon/dateparse"
	"github.com/ppowo/feedlet/internal/models"
)

// Compiled regexes for performance (compile once at package load time)
var (
	dateRegex       = regexp.MustCompile(`(?i)\b(jan(uary)?|feb(ruary)?|mar(ch)?|apr(il)?|may|june?|july?|aug(ust)?|sep(tember)?|oct(ober)?|nov(ember)?|dec(ember)?)\b|20\d{2}`)
	promotionRegex  = regexp.MustCompile(`(?i)\b(WWE|AEW|NJPW|ROH|CMLL|IMPACT|TNA|NEW JAPAN|ALL ELITE|RING OF HONOR|CONSEJO MUNDIAL DE LUCHA LIBRE|AAA|MLW|NWA|GCW|PWG)\b`)
	ratingRegex     = regexp.MustCompile(`^(\d\.?\d{0,2}|\d+[¼½¾])$`)
	numericRegex    = regexp.MustCompile(`\d+\.?\d*`)
	citationRegex   = regexp.MustCompile(`\[\d+\]`)
)

// WikipediaSource implements the Source interface for Wikipedia pages
type WikipediaSource struct {
	name       string
	url        string
	limit      int
	ignoreDays bool
	minRating  float64
}

// NewWikipediaSource creates a new Wikipedia source
func NewWikipediaSource(name, url string, limit int) *WikipediaSource {
	return &WikipediaSource{
		name:       name,
		url:        url,
		limit:      limit,
		ignoreDays: true, // Wikipedia entries should ignore day filtering by default
		minRating:  5.0,  // Minimum rating to include (>= 5.0 gets all 5+ star matches)
	}
}

// Fetch retrieves match data from Wikipedia
func (w *WikipediaSource) Fetch(ctx context.Context) ([]models.Item, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", w.url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set User-Agent to avoid being blocked
	req.Header.Set("User-Agent", "Feedlet/1.0 (RSS aggregator)")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch wikipedia: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// Create a slice to hold all items with their row numbers for sorting
	type itemWithRowNum struct {
		item    models.Item
		rowNum  int
	}

	allItems := make([]itemWithRowNum, 0)

	// Try all wikitable tables and find the one with the highest numbered entries
	var bestTable *goquery.Selection
	var maxRowNum int

	doc.Find("table.wikitable, table.TBase").Each(func(i int, table *goquery.Selection) {
		table.Find("tr").Each(func(j int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() >= 1 {
				rowNumStr := strings.TrimSpace(cells.Eq(0).Text())
				if rowNum, err := strconv.Atoi(rowNumStr); err == nil {
					if rowNum > maxRowNum {
						maxRowNum = rowNum
						bestTable = table
					}
				}
			}
		})
	})

	// If we found a table with high-numbered entries, use it
	if bestTable != nil {
		// Track the current values for merged cells
		var currentDate, currentPromotion, currentEvent, currentRating string

		bestTable.Find("tr").Each(func(j int, row *goquery.Selection) {
			cells := row.Find("td")
			if cells.Length() < 2 {
				return
			}

			// Get the first cell text which should contain the unique numbers
			firstCellText := strings.TrimSpace(cells.Eq(0).Text())

			// Extract the main row number (first number in the cell)
			parts := strings.Fields(firstCellText)
			if len(parts) == 0 {
				return
			}

			rowNum, err := strconv.Atoi(parts[0])
			if err != nil {
				return // Skip if not a valid number
			}

			// Extract data from cells, handling merged cells
			colCount := cells.Length()

			// Reset variables for this row
			var dateText, matchText, promotionText, eventText, ratingText string

			// Process cells and identify their content type
			// Cell 0: row number, Cell 1: secondary number (ignore both)
			for i := 2; i < colCount; i++ {
				cellText := strings.TrimSpace(cells.Eq(i).Text())
				if cellText == "" {
					continue
				}

				// Identify and assign cell content based on pattern matching
				switch {
				case isLikelyDate(cellText):
					dateText = cellText
					currentDate = cellText
				case isLikelyRating(cellText):
					ratingText = cellText
					currentRating = cellText
				case isLikelyPromotion(cellText):
					promotionText = cellText
					currentPromotion = cellText
				case matchText == "":
					// First unidentified cell is likely the match
					matchText = cellText
				case eventText == "":
					// Second unidentified cell is likely the event
					eventText = cellText
					currentEvent = cellText
				}
			}

			// Fill in missing fields from merged cells
			if dateText == "" { dateText = currentDate }
			if promotionText == "" { promotionText = currentPromotion }
			if eventText == "" { eventText = currentEvent }
			if ratingText == "" { ratingText = currentRating }

			// Parse rating to filter for entries above minimum
			ratingValue := parseRating(ratingText)

			// Skip if we don't have essential data OR rating is below minimum
			if matchText == "" || matchText == firstCellText || ratingValue < w.minRating {
				return
			}


			// Parse date - handle various formats
			published := parseWikipediaDate(dateText)

			// Create description combining event and promotion info
			description := fmt.Sprintf("★%.2f", ratingValue)
			if eventText != "" {
				description = fmt.Sprintf("%s | %s", description, eventText)
			}
			if promotionText != "" && promotionText != eventText {
				description = fmt.Sprintf("%s | %s", description, promotionText)
			}
			if len(description) > 200 {
				description = description[:197] + "..."
			}

			item := models.Item{
				Title:       matchText,
				Link:        w.url,
				Description: description,
				Content:     fmt.Sprintf("Match: %s\nEvent: %s\nPromotion: %s\nDate: %s\nRating: ★%.2f", matchText, eventText, promotionText, dateText, ratingValue),
				Author:      "Dave Meltzer (Wikipedia)",
				Published:   published,
				SourceName:  w.name,
				SourceType:  "wikipedia",
				IgnoreDays:  w.ignoreDays,
			}

			allItems = append(allItems, itemWithRowNum{
				item:   item,
				rowNum: rowNum,
			})
		})
	}

	// Sort by row number descending (highest numbers first = latest entries)
	sort.Slice(allItems, func(i, j int) bool {
		return allItems[i].rowNum > allItems[j].rowNum
	})

	// Take only the limit number of items (latest entries)
	items := make([]models.Item, 0)
	for i, itemWithRow := range allItems {
		if w.limit > 0 && i >= w.limit {
			break
		}
		items = append(items, itemWithRow.item)
	}

	return items, nil
}

// parseWikipediaDate parses various date formats found on Wikipedia
func parseWikipediaDate(dateStr string) time.Time {
	// Remove citation markers like [1], [2], etc.
	cleanDate := citationRegex.ReplaceAllString(dateStr, "")
	cleanDate = strings.TrimSpace(cleanDate)

	if cleanDate == "" {
		return time.Now()
	}

	if t, err := dateparse.ParseAny(cleanDate); err == nil {
		return t
	}
	return time.Now()
}

// parseRating converts rating text to a float value
func parseRating(ratingText string) float64 {
	if ratingText == "" {
		return 0.0
	}

	// Try direct float parsing first
	if parsed, err := strconv.ParseFloat(ratingText, 64); err == nil {
		return parsed
	}

	// Handle special fraction symbols
	fractionMap := map[string]float64{
		"5¼": 5.25, "5.25": 5.25,
		"5½": 5.5, "5.5": 5.5,
		"5¾": 5.75, "5.75": 5.75,
		"6¼": 6.25, "6.25": 6.25,
		"6½": 6.5, "6.5": 6.5,
		"6¾": 6.75, "6.75": 6.75,
	}
	if val, ok := fractionMap[ratingText]; ok {
		return val
	}

	// Try to extract numeric part
	if matches := numericRegex.FindStringSubmatch(ratingText); len(matches) > 0 {
		if parsed, err := strconv.ParseFloat(matches[0], 64); err == nil {
			return parsed
		}
	}

	return 0.0
}

// Helper functions to identify cell content types
func isLikelyDate(text string) bool {
	return dateRegex.MatchString(text)
}

func isLikelyPromotion(text string) bool {
	return promotionRegex.MatchString(text)
}

func isLikelyRating(text string) bool {
	return ratingRegex.MatchString(text)
}

// Name returns the source name
func (w *WikipediaSource) Name() string {
	return w.name
}

// Type returns the source type
func (w *WikipediaSource) Type() string {
	return "wikipedia"
}