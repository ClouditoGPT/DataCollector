# DataCollector

A generic HTML web crawler with multi-source support.

## Architecture

### `internal/crawler/`
- `collector.go` - Main collector with `SourceFetcher` interface
- `html_client.go` - Generic `Client` for any HTML page
- `state.go` - Runtime state tracking (visited, queue, errors)
- `logger/` - Thread-safe logging

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
- Auto language detection (Persian/Arabic vs Latin)
- Logging with state
- Persistent queue/visited storage
- Graceful shutdown (SIGINT/SIGTERM)