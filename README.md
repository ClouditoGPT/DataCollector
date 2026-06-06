# DataCollector

A generic multi-source web crawler with support for Wikipedia and HTML websites.

## Architecture

### Generic Crawler (`internal/crawler/`)
- `collector.go` - Main collector with `SourceFetcher` interface
- `queue.go` - Thread-safe URL queue management
- `visited.go` - Thread-safe visited URL tracking
- `queue_store.go` / `visited_store.go` - Persistent JSON storage
- `ratelimit.go` - Rate limiting for requests
- `state.go` - Runtime state tracking (visited count, queue size, errors)

### Logger (`internal/logger/`)
- `logger.go` - Thread-safe logging with Info/Error/Debug levels

### Sources

Each source implements `SourceFetcher` interface:
- `internal/sources/wikipedia/client.go` - Wikipedia API client
- `internal/sources/htmlcrawler/client.go` - HTML website client

## Configuration (`configs/sources.json`)

```json
{
  "sources": [
    {
      "name": "wikipedia",
      "url": "https://fa.wikipedia.org/w/api.php",
      "seeds": ["ایران", "تهران"],
      "workers": 5,
      "rate_delay_ms": 500
    },
    {
      "name": "html",
      "url": "https://example.com",
      "seeds": ["https://example.com"],
      "workers": 3,
      "rate_delay_ms": 1000
    }
  ]
}
```

## Features

- **Auto language detection**: Detects Persian/Farsi (`fa`) or English (`en`) from content
- **State tracking**: Tracks visited pages, queue size, and errors per source
- **Persistent state**: Resumes from saved queue/visited on restart
- **Concurrent workers**: Configurable worker count per source