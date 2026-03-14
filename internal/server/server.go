package server

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"sort"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/ppowo/feedlet/internal/aggregator"
	"github.com/ppowo/feedlet/internal/fetcher"
	"github.com/ppowo/feedlet/internal/knowledge"
	"github.com/ppowo/feedlet/internal/models"
)

// Server represents the HTTP server.
type Server struct {
	fetcher       *fetcher.Fetcher
	tmpl          *template.Template
	port          int
	sourceConfigs []models.SourceConfig
	sourceLimits  map[string]int
	sourceDays    map[string]int
	httpServer    *http.Server
}

// New creates a new server from embedded template content.
func New(f *fetcher.Fetcher, templateContent string, port int, sourceConfigs []models.SourceConfig, sourceLimits map[string]int, sourceDays map[string]int) (*Server, error) {
	funcMap := template.FuncMap{
		"formatTime":    func(t time.Time) string { return t.Format("Jan 2, 2006 3:04 PM") },
		"formatTimeAgo": humanize.Time,
	}

	tmpl, err := template.New("index.html").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	s := &Server{
		fetcher:       f,
		tmpl:          tmpl,
		port:          port,
		sourceConfigs: append([]models.SourceConfig(nil), sourceConfigs...),
		sourceLimits:  sourceLimits,
		sourceDays:    sourceDays,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.handleIndex)
	mux.HandleFunc("/events", s.handleSSE)

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return s, nil
}

func (s *Server) Start() error {
	log.Printf("Starting Feedlet")
	log.Printf("Dashboard: http://localhost:%d", s.port)

	for _, url := range localAccessURLs(s.port) {
		log.Printf("Dashboard (LAN): %s", url)
	}

	log.Printf("Press Ctrl+C to stop")

	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down HTTP server...")
	return s.httpServer.Shutdown(ctx)
}

func localAccessURLs(port int) []string {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil
	}

	urls := make([]string, 0)
	seen := make(map[string]bool)

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			url := fmt.Sprintf("http://%s:%d", ip.String(), port)
			if !seen[url] {
				urls = append(urls, url)
				seen[url] = true
			}
		}
	}

	sort.Strings(urls)
	return urls
}

// handleSSE serves server-sent events for feed updates.
func (s *Server) handleSSE(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	updateCh, err := s.fetcher.Subscribe()
	if err != nil {
		log.Printf("Failed to subscribe: %v", err)
		http.Error(w, "Server at capacity", http.StatusServiceUnavailable)
		return
	}
	defer s.fetcher.Unsubscribe(updateCh)

	fmt.Fprintf(w, "data: ping\n\n")
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

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

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	feed := s.fetcher.GetFeed()
	filteredAgg := aggregator.Process(feed).FilterBySourceDays(s.sourceDays)
	filteredGrouped := filteredAgg.GroupBySource()
	grouped := filteredAgg.LimitPerSource(s.sourceLimits).GroupBySource()
	allGrouped := aggregator.Process(feed).GroupBySource()

	type Source struct {
		Name                string
		Items               []any
		HasItems            bool
		IgnoreDays          bool
		NSFW                bool
		IsChronological     bool
		Days                int
		HiddenCount         int
		NewestItemAge       time.Time
		Error               string
		Stale               bool
		LastAttemptAt       time.Time
		LastSuccessAt       time.Time
		ConsecutiveFailures int
		HasEverSucceeded    bool
		HasHistoricalItems  bool
		IsWaiting           bool
		StatusText          string
		EmptyText           string
		ShowStaleBanner     bool
		ShowErrorPanel      bool
		Order               int
	}

	getDays := func(name string) int {
		days := s.sourceDays[name]
		if days == 0 {
			return 2
		}
		return days
	}

	applyState := func(dst *Source, name string) {
		if state, ok := feed.SourceStates[name]; ok {
			dst.Error = state.LastError
			dst.Stale = state.Stale
			dst.LastAttemptAt = state.LastAttemptAt
			dst.LastSuccessAt = state.LastSuccessAt
			dst.ConsecutiveFailures = state.ConsecutiveFailures
			dst.HasEverSucceeded = !state.LastSuccessAt.IsZero()
		}
		if dst.Error == "" {
			if errMsg, ok := feed.Errors[name]; ok {
				dst.Error = errMsg
				dst.Stale = true
			}
		}
	}

	sourceByName := make(map[string]*Source, len(s.sourceConfigs))
	ordered := make([]*Source, 0, len(s.sourceConfigs))

	for i, cfg := range s.sourceConfigs {
		src := &Source{
			Name:            cfg.Name,
			Items:           []any{},
			IgnoreDays:      cfg.IgnoreDays,
			NSFW:            cfg.NSFW,
			IsChronological: cfg.IsChronological,
			Days:            getDays(cfg.Name),
			Order:           i,
		}
		applyState(src, cfg.Name)
		sourceByName[cfg.Name] = src
		ordered = append(ordered, src)
	}

	ensureSource := func(name string) *Source {
		if src, ok := sourceByName[name]; ok {
			return src
		}

		src := &Source{
			Name:  name,
			Items: []any{},
			Days:  getDays(name),
			Order: len(ordered),
		}
		applyState(src, name)
		sourceByName[name] = src
		ordered = append(ordered, src)
		return src
	}

	for name, items := range grouped {
		src := ensureSource(name)
		src.HiddenCount = len(filteredGrouped[name]) - len(items)
		src.Items = make([]any, len(items))
		src.HasItems = len(items) > 0

		for i, item := range items {
			src.Items[i] = item
			if item.IgnoreDays {
				src.IgnoreDays = true
			}
			if item.NSFW {
				src.NSFW = true
			}
			if item.IsChronological {
				src.IsChronological = true
			}
			if item.Published.After(src.NewestItemAge) {
				src.NewestItemAge = item.Published
			}
		}
	}

	for name := range allGrouped {
		src := ensureSource(name)
		src.HasHistoricalItems = len(allGrouped[name]) > 0
		if src.NewestItemAge.IsZero() {
			for _, item := range allGrouped[name] {
				if item.Published.After(src.NewestItemAge) {
					src.NewestItemAge = item.Published
				}
			}
		}
	}
	for name := range feed.SourceStates {
		ensureSource(name)
	}
	for name := range feed.Errors {
		ensureSource(name)
	}

	for _, src := range ordered {
		src.IsWaiting = !src.HasItems && !src.HasEverSucceeded && src.Error == "" && src.LastAttemptAt.IsZero()
		src.ShowStaleBanner = src.HasItems && src.Stale
		src.ShowErrorPanel = !src.HasItems && src.Error != ""

		switch {
		case src.Stale && src.HasEverSucceeded:
			src.StatusText = humanize.Time(src.LastSuccessAt)
		case src.Stale && !src.LastAttemptAt.IsZero():
			src.StatusText = "refresh failed"
		case src.HasEverSucceeded:
			src.StatusText = humanize.Time(src.LastSuccessAt)
		case src.IsWaiting:
			src.StatusText = "waiting"
		case !src.LastAttemptAt.IsZero() && src.Error == "":
			src.StatusText = "fetching"
		default:
			src.StatusText = "waiting"
		}

		if !src.HasItems && !src.ShowErrorPanel {
			switch {
			case src.IsWaiting:
				src.EmptyText = "Waiting for first fetch..."
			case src.HasEverSucceeded || src.HasHistoricalItems || !src.LastAttemptAt.IsZero():
				src.EmptyText = "No recent items"
			default:
				src.EmptyText = "Waiting for first fetch..."
			}
		}
	}

	sources := make([]Source, 0, len(ordered))
	for _, src := range ordered {
		sources = append(sources, *src)
	}

	sort.SliceStable(sources, func(i, j int) bool {
		if sources[i].NewestItemAge.Equal(sources[j].NewestItemAge) {
			return sources[i].Order < sources[j].Order
		}
		if sources[i].NewestItemAge.IsZero() {
			return false
		}
		if sources[j].NewestItemAge.IsZero() {
			return true
		}
		return sources[i].NewestItemAge.After(sources[j].NewestItemAge)
	})

	knowledgeBit := knowledge.GetRandomBit()

	data := struct {
		Sources      []Source
		UpdatedAt    time.Time
		Title        string
		KnowledgeBit knowledge.KnowledgeBit
	}{
		Sources:      sources,
		UpdatedAt:    feed.UpdatedAt,
		Title:        "Feedlet",
		KnowledgeBit: knowledgeBit,
	}

	if err := s.tmpl.Execute(w, data); err != nil {
		log.Printf("Error rendering template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
