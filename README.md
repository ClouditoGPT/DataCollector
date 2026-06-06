# DataCollector

A generic HTML web crawler.

## Architecture

### `internal/crawler/`
- `collector.go` - Main collector with `SourceFetcher` interface, auto language detection, state tracking, raw page saving
- `html_client.go` - Generic `Client` for any HTML page
- `queue.go` - Thread-safe URL queue
- `visited.go` - Thread-safe visited set
- `queue_store.go` / `visited_store.go` - Persistent JSON storage
- `ratelimit.go` - Rate limiting
- `state.go` - Runtime state tracking
- `logger.go` - Logging

## Configuration

`configs/sources.json`:
```json
{
  "sources": [
    {
      "name": "html",
      "url": "https://en.wikipedia.org",
      "seeds": ["https://en.wikipedia.org/wiki/Iran"],
      "workers": 5,
      "rate_delay_ms": 500
    }
  ]
}
```

## Usage

```bash
go run ./cmd/collector/main.go
```

## Features
- Auto language detection
- Raw pages saved to `./data/raw/{source}/{id}.json`
- State persistence
- Graceful shutdown (SIGINT/SIGTERM)