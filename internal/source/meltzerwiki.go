package source

import (
	"context"
	"fmt"
	"net/http"
	neturl "net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/ppowo/feedlet/internal/models"
	"github.com/ppowo/feedlet/internal/source/httpclient"
)

const (
	meltzerWikiPageURL   = "https://en.wikipedia.org/wiki/List_of_professional_wrestling_matches_rated_5_or_more_stars_by_Dave_Meltzer"
	meltzerWikiRenderURL = "https://en.wikipedia.org/w/index.php?title=List_of_professional_wrestling_matches_rated_5_or_more_stars_by_Dave_Meltzer&action=render"
)

var meltzerWikiBaseURL, _ = neturl.Parse("https://en.wikipedia.org/")

// CurrentMeltzerWikiSourceName returns the display name for the Wikipedia-based Meltzer source.
func CurrentMeltzerWikiSourceName() string {
	return "Latest Meltzer 5★+ Matches"
}

// CurrentMeltzerWikiHomeURL returns the canonical article URL for the source.
func CurrentMeltzerWikiHomeURL() string {
	return meltzerWikiPageURL
}

// MeltzerWikiSource fetches the latest Dave Meltzer 5+ star matches from Wikipedia.
type MeltzerWikiSource struct {
	name  string
	limit int
}

// NewMeltzerWikiSource creates a new MeltzerWiki source.
func NewMeltzerWikiSource(name string, limit int) *MeltzerWikiSource {
	return &MeltzerWikiSource{
		name:  name,
		limit: limit,
	}
}

type meltzerWikiCell struct {
	text string
	link string
}

type meltzerWikiPendingCell struct {
	cell      meltzerWikiCell
	remaining int
}

func resolveMeltzerWikiURL(href string) string {
	ref, err := neturl.Parse(strings.TrimSpace(href))
	if err != nil || ref.String() == "" {
		return ""
	}
	if meltzerWikiBaseURL == nil {
		return ref.String()
	}
	return meltzerWikiBaseURL.ResolveReference(ref).String()
}

func normalizeWhitespace(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func cleanMeltzerWikiCellText(cell *goquery.Selection) string {
	html, err := cell.Html()
	if err != nil {
		return normalizeWhitespace(cell.Text())
	}

	html = strings.NewReplacer("<br>", " / ", "<br/>", " / ", "<br />", " / ").Replace(html)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<div>" + html + "</div>"))
	if err != nil {
		return normalizeWhitespace(cell.Text())
	}
	doc.Find("sup.reference").Remove()
	return normalizeWhitespace(doc.Text())
}

func extractMeltzerWikiCell(cell *goquery.Selection) meltzerWikiCell {
	return meltzerWikiCell{
		text: cleanMeltzerWikiCellText(cell),
		link: resolveMeltzerWikiURL(strings.TrimSpace(cell.Find("a[href]").First().AttrOr("href", ""))),
	}
}

func consumeMeltzerWikiPendingCells(cells []meltzerWikiCell, pending map[int]meltzerWikiPendingCell, col *int) []meltzerWikiCell {
	for {
		span, ok := pending[*col]
		if !ok {
			return cells
		}

		cells = append(cells, span.cell)
		if span.remaining <= 1 {
			delete(pending, *col)
		} else {
			span.remaining--
			pending[*col] = span
		}
		*col++
	}
}

func parseMeltzerWikiTable(table *goquery.Selection) ([][]meltzerWikiCell, error) {
	rows := make([][]meltzerWikiCell, 0, 64)
	pending := make(map[int]meltzerWikiPendingCell)

	table.Find("tr").Each(func(_ int, row *goquery.Selection) {
		if row.Find("td").Length() == 0 || row.Find("th").Length() > 0 {
			return
		}

		cells := make([]meltzerWikiCell, 0, 8)
		col := 0
		row.ChildrenFiltered("td").Each(func(_ int, cell *goquery.Selection) {
			cells = consumeMeltzerWikiPendingCells(cells, pending, &col)

			parsed := extractMeltzerWikiCell(cell)
			cells = append(cells, parsed)

			if rowspan, err := strconv.Atoi(strings.TrimSpace(cell.AttrOr("rowspan", ""))); err == nil && rowspan > 1 {
				pending[col] = meltzerWikiPendingCell{cell: parsed, remaining: rowspan - 1}
			}
			col++
		})

		cells = consumeMeltzerWikiPendingCells(cells, pending, &col)
		if len(cells) > 0 {
			rows = append(rows, cells)
		}
	})

	if len(rows) == 0 {
		return nil, fmt.Errorf("meltzerwiki: no data rows found in matches table")
	}

	return rows, nil
}

func isMeltzerWikiMatchesTable(table *goquery.Selection) bool {
	headerText := normalizeWhitespace(table.Find("tr").First().Text())
	return strings.Contains(headerText, "Date") &&
		strings.Contains(headerText, "Match") &&
		strings.Contains(headerText, "Promotion") &&
		strings.Contains(headerText, "Event") &&
		strings.Contains(headerText, "Rating")
}

func (m *MeltzerWikiSource) fetchDocument(ctx context.Context) (*goquery.Document, error) {
	req, err := retryablehttp.NewRequestWithContext(ctx, http.MethodGet, meltzerWikiRenderURL, nil)
	if err != nil {
		return nil, fmt.Errorf("meltzerwiki: failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", httpclient.RandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := httpclient.GetClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("meltzerwiki: failed to fetch listing: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("meltzerwiki: unexpected status %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("meltzerwiki: failed to parse HTML: %w", err)
	}

	return doc, nil
}

func (m *MeltzerWikiSource) parseMatchRow(row []meltzerWikiCell) (models.Item, bool) {
	if len(row) < 7 {
		return models.Item{}, false
	}

	dateStr := row[2].text
	matchTitle := row[3].text
	promotion := row[4].text
	event := row[5].text
	rating := row[6].text
	if dateStr == "" || matchTitle == "" || promotion == "" || event == "" || rating == "" {
		return models.Item{}, false
	}

	published, err := time.Parse("January 2, 2006", dateStr)
	if err != nil {
		return models.Item{}, false
	}

	link := row[5].link
	if link == "" {
		link = CurrentMeltzerWikiHomeURL()
	}

	return models.Item{
		Title:       matchTitle,
		Link:        link,
		Description: buildMeltzerWikiDescription(rating, promotion, event, dateStr),
		Content:     buildMeltzerWikiContent(matchTitle, promotion, event, rating, dateStr),
		Author:      "Wikipedia",
		Published:   published,
		SourceName:  m.name,
		SourceType:  "meltzerwiki",
		IgnoreDays:  true,
	}, true
}

func (m *MeltzerWikiSource) parseDocument(doc *goquery.Document) ([]models.Item, error) {
	items := make([]models.Item, 0, 256)
	matchTables := 0
	var parseErr error

	doc.Find("table.wikitable.sortable").Each(func(_ int, table *goquery.Selection) {
		if parseErr != nil || !isMeltzerWikiMatchesTable(table) {
			return
		}

		matchTables++
		rows, err := parseMeltzerWikiTable(table)
		if err != nil {
			parseErr = err
			return
		}

		for _, row := range rows {
			if item, ok := m.parseMatchRow(row); ok {
				items = append(items, item)
			}
		}
	})

	if parseErr != nil {
		return nil, parseErr
	}
	if matchTables == 0 {
		return nil, fmt.Errorf("meltzerwiki: matches table not found")
	}
	if len(items) == 0 {
		return nil, fmt.Errorf("meltzerwiki: no match rows parsed")
	}

	return items, nil
}

// Fetch retrieves the latest Meltzer 5+ star matches from Wikipedia.
func (m *MeltzerWikiSource) Fetch(ctx context.Context) ([]models.Item, error) {
	doc, err := m.fetchDocument(ctx)
	if err != nil {
		return nil, err
	}

	items, err := m.parseDocument(doc)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(items, func(i, j int) bool {
		return items[i].Published.After(items[j].Published)
	})

	if m.limit > 0 && len(items) > m.limit {
		items = items[:m.limit]
	}

	return items, nil
}

// Name returns the source name.
func (m *MeltzerWikiSource) Name() string { return m.name }

// Type returns the source type.
func (m *MeltzerWikiSource) Type() string { return "meltzerwiki" }

func buildMeltzerWikiDescription(rating, promotion, event, dateStr string) string {
	parts := []string{fmt.Sprintf("★%s", rating)}
	if promotion != "" {
		parts = append(parts, promotion)
	}
	if event != "" {
		parts = append(parts, event)
	}
	if dateStr != "" {
		parts = append(parts, dateStr)
	}
	return strings.Join(parts, " | ")
}

func buildMeltzerWikiContent(matchTitle, promotion, event, rating, dateStr string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "Match: %s\n", matchTitle)
	fmt.Fprintf(&b, "Promotion: %s\n", promotion)
	fmt.Fprintf(&b, "Event: %s\n", event)
	fmt.Fprintf(&b, "Rating: ★%s\n", rating)
	fmt.Fprintf(&b, "Date: %s\n", dateStr)
	return strings.TrimRight(b.String(), "\n")
}
