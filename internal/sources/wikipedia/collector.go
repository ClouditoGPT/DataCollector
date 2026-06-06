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

	topics := []string{
		"ایران",
		"تهران",
		"برنامه‌نویسی",
		"هوش مصنوعی",
		"زبان گو",
	}

	ch := make(chan models.Document)

	go func() {
		for _, topic := range topics {

			title, text, err := FetchArticle(topic)

			if err != nil {
				continue
			}

			doc := models.Document{
				ID:       uuid.NewString(),
				Source:   "wikipedia",
				Type:     models.ArticleDocument,
				Language: "fa",
				Title:    title,
				Content:  text,
			}

			select {
			case <-ctx.Done():
				return

			case ch <- doc:
			}
		}
	}()

	return ch, nil
}
