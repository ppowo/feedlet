package source

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

const (
	cagematchBaseURL   = "https://www.cagematch.net/"
	cagematchMinRating = 8.75
	cagematchMinVotes  = 500
)

// CurrentCagematchSourceName returns the display name derived from the source thresholds.
func CurrentCagematchSourceName() string {
	return fmt.Sprintf("Top Matches (★%.1f+ / %d+ votes)", cagematchMinRating, cagematchMinVotes)
}

// CagematchSource fetches the most recent top-rated matches from Cagematch.
// It queries the Matchguide filtered to the current year, sorted by date
// descending, keeping only matches above the minimum rating and vote thresholds.
type CagematchSource struct {
	name  string
	limit int
}

// NewCagematchSource creates a new Cagematch source.
func NewCagematchSource(name string, limit int) *CagematchSource {
	return &CagematchSource{
		name:  name,
		limit: limit,
	}
}

// CurrentCagematchHomeURL returns the current-year Matchguide URL used by the source.
func CurrentCagematchHomeURL() string {
	return cagematchHomeURLForYear(time.Now().Year())
}

func cagematchHomeURLForYear(year int) string {
	return fmt.Sprintf(
		"%s?id=112&view=list&year=%d&minRating=%.1f&minVotes=%d&sortby=colDate&sorttype=DESC",
		cagematchBaseURL,
		year,
		cagematchMinRating,
		cagematchMinVotes,
	)
}

// buildURL constructs the Matchguide URL for the current year.
func (c *CagematchSource) buildURL() string {
	return CurrentCagematchHomeURL()
}

func (c *CagematchSource) Fetch(ctx context.Context) ([]models.Item, error) {
	fetchURL := c.buildURL()

	req, err := retryablehttp.NewRequestWithContext(ctx, "GET", fetchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("cagematch: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", httpclient.RandomUserAgent())

	resp, err := httpclient.GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("cagematch: failed to fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("cagematch: unexpected status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("cagematch: failed to parse HTML: %w", err)
	}

	items := make([]models.Item, 0)

	doc.Find("tr").Each(func(_ int, row *goquery.Selection) {
		if c.limit > 0 && len(items) >= c.limit {
			return
		}

		cells := row.Find("td")
		if cells.Length() < 8 {
			return // skip header rows and short rows
		}

		// Col 0: row number — skip if not numeric (shouldn't happen, but guards stray rows)
		if _, err := strconv.Atoi(strings.TrimSpace(cells.Eq(0).Text())); err != nil {
			return
		}

		// Col 1: date in DD.MM.YYYY format
		dateStr := strings.TrimSpace(cells.Eq(1).Text())
		published, err := time.Parse("02.01.2006", dateStr)
		if err != nil {
			return
		}

		// Col 2: promotion — collect alt text from all promotion images
		var promotionParts []string
		cells.Eq(2).Find("img").Each(func(_ int, img *goquery.Selection) {
			if alt := strings.TrimSpace(img.AttrOr("alt", "")); alt != "" {
				promotionParts = append(promotionParts, alt)
			}
		})
		promotion := strings.Join(promotionParts, " / ")

		// Col 3: match fixture — title from link text, href for the detail page
		matchCell := cells.Eq(3)
		matchTitle := strings.TrimSpace(matchCell.Text())
		matchHref := strings.TrimSpace(matchCell.Find("a").First().AttrOr("href", ""))
		matchLink := cagematchBaseURL
		if strings.HasPrefix(matchHref, "http") {
			matchLink = matchHref
		} else if matchHref != "" {
			// relative href like "?id=111&nr=XXXXX"
			matchLink = cagematchBaseURL + matchHref
		}

		// Col 4: WON (Meltzer) star rating — may be empty
		wonRating := strings.TrimSpace(cells.Eq(4).Text())

		// Col 5: match type
		matchType := strings.TrimSpace(cells.Eq(5).Text())

		// Col 6: Cagematch community rating — validate against minimum
		ratingStr := strings.TrimSpace(cells.Eq(6).Text())
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil || rating < cagematchMinRating {
			return
		}

		// Col 7: vote count — validate against minimum
		votesStr := strings.TrimSpace(cells.Eq(7).Text())
		votes, err := strconv.Atoi(votesStr)
		if err != nil || votes < cagematchMinVotes {
			return
		}

		description := buildCagematchDescription(rating, votes, matchType, wonRating, promotion, dateStr)
		content := buildCagematchContent(matchTitle, promotion, dateStr, rating, votes, matchType, wonRating)

		items = append(items, models.Item{
			Title:       matchTitle,
			Link:        matchLink,
			Description: description,
			Content:     content,
			Author:      "Cagematch",
			Published:   published,
			SourceName:  c.name,
			SourceType:  "cagematch",
			IgnoreDays:  true,
		})
	})

	return items, nil
}

func (c *CagematchSource) Name() string { return c.name }
func (c *CagematchSource) Type() string { return "cagematch" }

func buildCagematchDescription(rating float64, votes int, matchType, wonRating, promotion, dateStr string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "★%.2f (%d votes)", rating, votes)
	if matchType != "" {
		fmt.Fprintf(&b, " | %s", matchType)
	}
	if wonRating != "" {
		fmt.Fprintf(&b, " | WON: %s", wonRating)
	}
	if promotion != "" {
		fmt.Fprintf(&b, " | %s", promotion)
	}
	fmt.Fprintf(&b, " | %s", dateStr)
	return b.String()
}

func buildCagematchContent(matchTitle, promotion, dateStr string, rating float64, votes int, matchType, wonRating string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Match: %s\n", matchTitle)
	fmt.Fprintf(&b, "Promotion: %s\n", promotion)
	fmt.Fprintf(&b, "Date: %s\n", dateStr)
	fmt.Fprintf(&b, "Rating: ★%.2f\n", rating)
	fmt.Fprintf(&b, "Votes: %d\n", votes)
	if matchType != "" {
		fmt.Fprintf(&b, "Type: %s\n", matchType)
	}
	if wonRating != "" {
		fmt.Fprintf(&b, "WON: %s\n", wonRating)
	}
	return strings.TrimRight(b.String(), "\n")
}
