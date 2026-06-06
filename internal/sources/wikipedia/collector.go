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

	topics := []string{"ایران", "تهران", "برنامه‌نویسی", "هوش مصنوعی", "زبان گو"}

	go func() {
		defer close(ch)

		for _, topic := range topics {

			title, text, err := FetchArticle(topic)
			if err != nil {
				continue
			}

			if text == "" {
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

				select {
				case <-ctx.Done():
					return
				default:
					// DO NOT create fake documents yet
					continue
				}
			}
		}
	}()

	return ch, nil
}
