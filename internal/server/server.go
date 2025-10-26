package server

import (
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/ppowo/feedlet/internal/aggregator"
	"github.com/ppowo/feedlet/internal/fetcher"
	"github.com/ppowo/feedlet/internal/knowledge"
)

// Server represents the HTTP server
type Server struct {
	fetcher      *fetcher.Fetcher
	tmpl         *template.Template
	port         int
	sourceLimits map[string]int
	sourceDays   map[string]int
}

// New creates a new server from embedded template content
func New(f *fetcher.Fetcher, templateContent string, port int, sourceLimits map[string]int, sourceDays map[string]int) (*Server, error) {
	// Custom template functions
	funcMap := template.FuncMap{
		"formatTime": func(t time.Time) string {
			return t.Format("Jan 2, 2006 3:04 PM")
		},
		"formatTimeAgo": func(t time.Time) string {
			duration := time.Since(t)
			switch {
			case duration < time.Minute:
				return "just now"
			case duration < time.Hour:
				minutes := int(duration.Minutes())
				if minutes == 1 {
					return "1 minute ago"
				}
				return fmt.Sprintf("%d minutes ago", minutes)
			case duration < 24*time.Hour:
				hours := int(duration.Hours())
				if hours == 1 {
					return "1 hour ago"
				}
				return fmt.Sprintf("%d hours ago", hours)
			case duration < 7*24*time.Hour:
				days := int(duration.Hours() / 24)
				if days == 1 {
					return "1 day ago"
				}
				return fmt.Sprintf("%d days ago", days)
			default:
				return t.Format("Jan 2, 2006")
			}
		},
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &Server{
		fetcher:      f,
		tmpl:         tmpl,
		port:         port,
		sourceLimits: sourceLimits,
		sourceDays:   sourceDays,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	http.HandleFunc("/", s.handleIndex)
	http.HandleFunc("/events", s.handleSSE)

	addr := fmt.Sprintf(":%d", s.port)
	log.Printf("Starting server on http://localhost%s", addr)
	return http.ListenAndServe(addr, nil)
}

// handleSSE serves server-sent events for feed updates
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	// Set headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Subscribe to feed updates
	updateCh := s.fetcher.Subscribe()
	defer s.fetcher.Unsubscribe(updateCh)

	// Send ping on connect
	fmt.Fprintf(w, "data: ping\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// Listen for updates or client disconnect
	for {
		select {
		case <-updateCh:
			fmt.Fprintf(w, "data: update\n\n")
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
		case <-r.Context().Done():
			return
		}
	}
}

// handleIndex serves the main page
func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	feed := s.fetcher.GetFeed()
	agg := aggregator.Process(feed)

	// Filter by per-source days
	agg = agg.FilterBySourceDays(s.sourceDays)

	// Apply per-source limits
	agg = agg.LimitPerSource(s.sourceLimits)

	// Group items by source
	grouped := agg.GroupBySource()

	// Convert to slice of sources for template
	// Include ALL enabled sources, even if they have no items after filtering
	type Source struct {
		Name          string
		Items         []any
		IgnoreDays    bool
		NSFW          bool
		Days          int
		NewestItemAge time.Time // For sorting by freshness
		Error         string    // Error message if fetch failed
	}

	// Get all enabled source names from config
	allSourceNames := make(map[string]bool)
	for _, item := range agg.Items {
		allSourceNames[item.SourceName] = true
	}

	sources := make([]Source, 0)
	for name, items := range grouped {
		itemsInterface := make([]any, len(items))
		ignoreDays := false
		nsfw := false
		var newestTime time.Time
		for i, item := range items {
			itemsInterface[i] = item
			if item.IgnoreDays {
				ignoreDays = true
			}
			if item.NSFW {
				nsfw = true
			}
			// Track the newest item's publish time
			if item.Published.After(newestTime) {
				newestTime = item.Published
			}
		}
		days := s.sourceDays[name]
		if days == 0 {
			days = 2 // Default
		}

		// Check if there's an error for this source
		errorMsg := ""
		if errMsg, hasError := feed.Errors[name]; hasError {
			errorMsg = errMsg
		}

		sources = append(sources, Source{
			Name:          name,
			Items:         itemsInterface,
			IgnoreDays:    ignoreDays,
			NSFW:          nsfw,
			Days:          days,
			NewestItemAge: newestTime,
			Error:         errorMsg,
		})
		delete(allSourceNames, name)
	}

	// Add sources with no items (but are configured and enabled)
	// We need to track which sources exist from the original feed before filtering
	feedBeforeFilter := s.fetcher.GetFeed()
	allSources := make(map[string]bool)
	sourceIgnoreDays := make(map[string]bool)
	sourceNSFW := make(map[string]bool)
	for _, item := range feedBeforeFilter.Items {
		allSources[item.SourceName] = true
		if item.IgnoreDays {
			sourceIgnoreDays[item.SourceName] = true
		}
		if item.NSFW {
			sourceNSFW[item.SourceName] = true
		}
	}

	for sourceName := range allSources {
		if _, exists := grouped[sourceName]; !exists {
			days := s.sourceDays[sourceName]
			if days == 0 {
				days = 2 // Default
			}

			// Check if there's an error for this source
			errorMsg := ""
			if errMsg, hasError := feed.Errors[sourceName]; hasError {
				errorMsg = errMsg
			}

			sources = append(sources, Source{
				Name:          sourceName,
				Items:         []any{},
				IgnoreDays:    sourceIgnoreDays[sourceName],
				NSFW:          sourceNSFW[sourceName],
				Days:          days,
				NewestItemAge: time.Time{}, // Zero time for sources with no items
				Error:         errorMsg,
			})
		}
	}

	// Also add sources that only have errors (not in allSources because no items ever fetched)
	for errorSourceName := range feed.Errors {
		found := false
		for _, src := range sources {
			if src.Name == errorSourceName {
				found = true
				break
			}
		}
		if !found {
			days := s.sourceDays[errorSourceName]
			if days == 0 {
				days = 2
			}
			sources = append(sources, Source{
				Name:          errorSourceName,
				Items:         []any{},
				IgnoreDays:    false,
				NSFW:          false,
				Days:          days,
				NewestItemAge: time.Time{},
				Error:         feed.Errors[errorSourceName],
			})
		}
	}

	// Sort sources by newest item time (most recent first)
	// Sources with no items (zero time) will appear last
	sort.Slice(sources, func(i, j int) bool {
		return sources[i].NewestItemAge.After(sources[j].NewestItemAge)
	})

	// Select a random knowledge bit
	knowledgeBits := knowledge.GetKnowledgeBits()
	randomBit := knowledgeBits[rand.Intn(len(knowledgeBits))]

	data := struct {
		Sources      []Source
		UpdatedAt    time.Time
		Title        string
		KnowledgeBit knowledge.KnowledgeBit
	}{
		Sources:      sources,
		UpdatedAt:    feed.UpdatedAt,
		Title:        "Feedlet",
		KnowledgeBit: randomBit,
	}

	if err := s.tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
