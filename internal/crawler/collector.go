package crawler

import (
	"DataCollector/internal/logger"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode"

	"github.com/google/uuid"
)

type DocumentType string

const (
	ArticleDocument     DocumentType = "article"
	ChatDocument        DocumentType = "chat"
	CodeDocument        DocumentType = "code"
	InstructionDocument DocumentType = "instruction"
	QADocument          DocumentType = "qa"
)

type Document struct {
	ID       string       `json:"id"`
	Source   string       `json:"source"`
	Type     DocumentType `json:"type"`
	Language string       `json:"language"`
	URL      string       `json:"url"`
	Title    string       `json:"title"`
	Content  any          `json:"content"`
}

type SourceFetcher interface {
	Name() string
	Fetch(url string) (title string, content string, links []string, err error)
}

type Collector struct {
	fetcher        SourceFetcher
	queueStorePath string
	visitedStorePath string
	rawStorePath   string
	seeds          []string
	rateDelay      time.Duration
	workers        int
	docType        DocumentType
	autoLangDetect bool
	state          *State
}

func NewCollector(fetcher SourceFetcher, opts ...func(*Collector)) *Collector {
	c := &Collector{
		fetcher:          fetcher,
		queueStorePath:   "./data/" + fetcher.Name() + "_queue.json",
		visitedStorePath: "./data/" + fetcher.Name() + "_visited.json",
		rawStorePath:     "./data/raw/" + fetcher.Name(),
		seeds:            []string{},
		rateDelay:        500 * time.Millisecond,
		workers:          5,
		docType:          ArticleDocument,
		autoLangDetect:   true,
		state:            GetState(fetcher.Name()),
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
		c.autoLangDetect = false
	}
}

func WithDocType(dt DocumentType) func(*Collector) {
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

func WithRawPath(path string) func(*Collector) {
	return func(c *Collector) {
		c.rawStorePath = path
	}
}

func (c *Collector) Name() string {
	return c.fetcher.Name()
}

func (c *Collector) GetState() *State {
	return c.state
}

func detectLanguage(content string) string {
	faCount := 0
	enCount := 0

	for _, r := range content {
		if unicode.Is(unicode.Arabic, r) || (r >= 0x0600 && r <= 0x06FF) {
			faCount++
		}
		if unicode.Is(unicode.Latin, r) {
			enCount++
		}
	}

	if faCount > enCount {
		return "fa"
	}
	return "en"
}

func saveRawPage(basePath, id string, page map[string]any) error {
	filename := filepath.Join(basePath, id+".json")
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (c *Collector) Collect(ctx context.Context) (<-chan Document, error) {
	ch := make(chan Document, 100)
	c.state.SetRunning(true)

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

	c.state.SetVisited(len(visitedMap))
	c.state.SetQueue(len(seed))

	limiter := NewRateLimiter(c.rateDelay)

	logger.Info("Starting crawler: source=%s, workers=%d, visited=%d, queue=%d", c.fetcher.Name(), c.workers, len(visitedMap), len(seed))

	var wg sync.WaitGroup

	worker := func() {
		defer wg.Done()
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
			c.state.SetVisited(c.state.GetVisited() + 1)

			logger.Info("Crawling: %s (queue=%d)", topic, len(queue.Snapshot()))

			limiter.Wait()

			title, text, links, err := c.fetcher.Fetch(topic)
			if err != nil || text == "" {
				if err != nil {
					logger.Error("Fetch failed for %s: %v", topic, err)
					c.state.IncErrors()
				}
				continue
			}

			lang := detectLanguage(text)

			id := uuid.NewString()
			doc := Document{
				ID:       id,
				Source:   c.fetcher.Name(),
				Type:     c.docType,
				Language: lang,
				URL:      topic,
				Title:    title,
				Content:  text,
			}

			ch <- doc

			logger.Info("Crawled: title=%s, language=%s, links=%d", title, lang, len(links))

			_ = saveRawPage(c.rawStorePath, id, map[string]any{
				"id":       id,
				"title":    title,
				"url":      topic,
				"language": lang,
				"content":  text,
			})

			for _, l := range links {
				if !visited.Has(l) {
					queue.Push(l)
					c.state.SetQueue(c.state.GetQueue() + 1)
				}
			}

			_ = queueStore.Save(queue.Snapshot())
			_ = visitedStore.Save(visited.Snapshot())
		}
	}

	wg.Add(c.workers)
	for i := 0; i < c.workers; i++ {
		go worker()
	}

	go func() {
		wg.Wait()
		close(ch)
		c.state.SetRunning(false)
		logger.Info("Crawler stopped: source=%s, visited=%d, errors=%d", c.fetcher.Name(), c.state.GetVisited(), c.state.GetErrors())
	}()

	return ch, nil
}