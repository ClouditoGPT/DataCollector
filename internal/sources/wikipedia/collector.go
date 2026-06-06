package wikipedia

import (
	"DataCollector/internal/crawler"
	"DataCollector/internal/models"
)

func New(baseURL string, seeds []string) *crawler.Collector {
	return crawler.NewCollector(
		NewClient(baseURL),
		crawler.WithSeeds(seeds),
		crawler.WithLanguage("fa"),
		crawler.WithDocType(models.ArticleDocument),
	)
}