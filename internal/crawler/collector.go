package crawler

import (
	"DataCollector/internal/models"
	"context"
	"time"

	"github.com/google/uuid"
)

type SourceFetcher interface {
	Name() string
	Fetch(url string) (title string, content string, links []string, err error)
}

type Collector struct {
	fetcher         SourceFetcher
	queueStorePath  string
	visitedStorePath string
	seeds           []string
	rateDelay       time.Duration
	workers         int
	language        string
	docType         models.DocumentType
}

func NewCollector(fetcher SourceFetcher, opts ...func(*Collector)) *Collector {
	c := &Collector{
		fetcher:         fetcher,
		queueStorePath:  "./data/" + fetcher.Name() + "_queue.json",
		visitedStorePath: "./data/" + fetcher.Name() + "_visited.json",
		seeds:           []string{},
		rateDelay:       500 * time.Millisecond,
		workers:         5,
		language:        "en",
		docType:         models.ArticleDocument,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func WithSeeds(seeds []string) func(*Collector) {
	return func(c *Collector) {
		c.seeds = seeds
	}
}

func WithRateDelay(d time.Duration) func(*Collector) {
	return func(c *Collector) {
		c.rateDelay = d
	}
}

func WithWorkers(n int) func(*Collector) {
	return func(c *Collector) {
		c.workers = n
	}
}

func WithLanguage(lang string) func(*Collector) {
	return func(c *Collector) {
		c.language = lang
	}
}

func WithDocType(dt models.DocumentType) func(*Collector) {
	return func(c *Collector) {
		c.docType = dt
	}
}

func WithQueuePath(path string) func(*Collector) {
	return func(c *Collector) {
		c.queueStorePath = path
	}
}

func WithVisitedPath(path string) func(*Collector) {
	return func(c *Collector) {
		c.visitedStorePath = path
	}
}

func (c *Collector) Name() string {
	return c.fetcher.Name()
}

func (c *Collector) Collect(ctx context.Context) (<-chan models.Document, error) {
	ch := make(chan models.Document, 100)

	queueStore := NewQueueStore(c.queueStorePath)
	visitedStore := NewVisitedStore(c.visitedStorePath)

	seed := c.seeds
	loadedQueue, err := queueStore.Load()
	if err == nil && len(loadedQueue) > 0 {
		seed = loadedQueue
	}

	queue := NewQueue(seed)
	visitedMap, _ := visitedStore.Load()
	visited := NewVisited()
	for k := range visitedMap {
		visited.Add(k)
	}

	limiter := NewRateLimiter(c.rateDelay)

	worker := func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			topic, ok := queue.Pop()
			if !ok {
				return
			}

			if visited.Has(topic) {
				continue
			}

			visited.Add(topic)

			limiter.Wait()

			title, text, links, err := c.fetcher.Fetch(topic)
			if err != nil || text == "" {
				continue
			}

			ch <- models.Document{
				ID:       uuid.NewString(),
				Source:   c.fetcher.Name(),
				Type:     c.docType,
				Language: c.language,
				URL:      topic,
				Title:    title,
				Content:  text,
			}

			for _, l := range links {
				if !visited.Has(l) {
					queue.Push(l)
				}
			}

			_ = queueStore.Save(queue.Snapshot())
			_ = visitedStore.Save(visited.Snapshot())
		}
	}

	for i := 0; i < c.workers; i++ {
		go worker()
	}

	go func() {
		<-ctx.Done()
		close(ch)
	}()

	return ch, nil
}