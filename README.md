# Feedlet

Minimal RSS feed aggregator with tiling layout and knowledge retention bar.

## Features

- Multiple source types (RSS, Reddit, Hacker News, custom scrapers)
- Auto-refresh via SSE
- Knowledge bar (random Java tidbits on each refresh)
- Embedded configuration

## Quick Start

```bash
# Development with auto-reload
mage dev

# Build for Linux x86_64
mage build

# Run
./target/feedlet
```

Server runs on `http://localhost:3737`

## Configuration

Edit `internal/config/config.go` and rebuild:

```go
{
    Name:           "r/programming",
    Type:           "reddit",
    URL:            "https://old.reddit.com/r/programming/top/.rss?t=month",
    Interval:       600,
    IntervalJitter: 120,
    Enabled:        true,
    Days:           7,
}
```

**Source types:** `rss`, `reddit`, `hnrss`, `wikipedia`

## Logging

Logs to stdout and OS log directory:

- macOS: `~/Library/Logs/feedlet/`
- Linux: `~/.local/state/feedlet/logs/`

Auto-rotates daily, keeps 3 days.

## Mage Commands

- `mage dev` - Run with auto-reload
- `mage build` - Build for Linux x86_64 to `target/`
- `mage clean` - Remove build artifacts
- `mage setup` - Install tools

## Knowledge Bar

Shows random Java concepts on each page load/refresh. Covers:

- OOP pillars (encapsulation, inheritance, polymorphism, abstraction)
- Collections, generics, exceptions
- Concurrency, streams, lambdas
- Design patterns, SOLID principles

Edit `internal/knowledge/knowledge.go` to customize.
