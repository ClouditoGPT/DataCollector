package wikipedia

import (
	"DataCollector/internal/models"
	"context"

	"github.com/google/uuid"
)

type Collector struct {
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) Name() string {
	return "wikipedia"
}

func (c *Collector) Collect(ctx context.Context) (<-chan models.Document, error) {

	ch := make(chan models.Document, 100)

	queueStore := NewQueueStore("./data/wikipedia_queue.json")
	visitedStore := NewVisitedStore("./data/wikipedia_visited.json")

	seed := []string{"ایران", "تهران", "هوش مصنوعی"}

	loadedQueue, err := queueStore.Load()
	if err == nil && len(loadedQueue) > 0 {
		seed = loadedQueue
	}

	queue := NewQueue(seed)

	visitedMap, _ := visitedStore.Load()
	visited := &Visited{set: visitedMap}

	go func() {
		defer close(ch)

		for {
			topic, ok := queue.Pop()
			if !ok {
				return
			}

			if visited.Has(topic) {
				continue
			}

			visited.Add(topic)

			title, text, err := FetchArticle(topic)
			if err != nil || text == "" {
				continue
			}

			ch <- models.Document{
				ID:       uuid.NewString(),
				Source:   "wikipedia",
				Type:     models.ArticleDocument,
				Language: "fa",
				Title:    title,
				Content:  text,
			}

			links, err := FetchLinks(topic)
			if err != nil {
				continue
			}

			for _, l := range links {
				if !visited.Has(l) {
					queue.Push(l)
				}
			}

			_ = queueStore.Save(queue.items)
			_ = visitedStore.Save(visited.set)
		}
	}()

	return ch, nil
}
