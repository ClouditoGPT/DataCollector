# DataCollector

A generic multi-source HTML web crawler.

## Architecture

### Generic Crawler (`internal/crawler/`)
- `html_client.go` - `HTMLClient` implementation of `SourceFetcher` for any HTML page
- `collector.go` - Main collector with auto language detection, rate limiting, state tracking
- `queue.go` - Thread-safe URL queue management
- `visited.go` - Thread-safe visited URL tracking
- `queue_store.go` / `visited_store.go` - Persistent JSON storage
- `ratelimit.go` - Rate limiting for requests
- `state.go` - Runtime state tracking per source
- `logger.go` - Thread-safe logging

## Usage

```bash
go run ./cmd/collector/main.go
```

Configure seeds in `cmd/collector/main.go`:
```go
seeds := []string{"https://en.wikipedia.org/wiki/Iran", "https://en.wikipedia.org/wiki/Tehran"}
```

## Features
- **Auto language detection**: Detects Persian (`fa`) vs English (`en`) from content
- **Logging**: All crawl events logged with topic, queue size, language
- **State persistence**: Queue and visited URLs saved to `./data/html_queue.json` and `./data/html_visited.json`
- **Concurrent workers**: Configurable worker count
- **Graceful shutdown**: Handles SIGINT/SIGTERM