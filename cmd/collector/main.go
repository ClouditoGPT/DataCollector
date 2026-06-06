package main

import (
	"DataCollector/internal/dedupe"
	"DataCollector/internal/pipeline"
	"DataCollector/internal/processors"
	"DataCollector/internal/storage"
	"context"
	"fmt"

	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	container.Provide(func() storage.Storage {
		return storage.NewJsonlStorage("./data")
	})
	container.Provide(func() (*dedupe.FileHashStore, error) {
		return dedupe.NewFileHashStore("./data/hashes.txt")
	})

	container.Provide(
		processors.NewDeduplicationProcessor,
	)
	container.Provide(
		processors.NewValidationProcessor,
	)

	container.Provide(
		func(
			store storage.Storage,
			validation *processors.ValidationProcessor,
			dedupe *processors.DeduplicationProcessor,
		) *pipeline.Pipeline {

			return pipeline.NewPipeline(
				store,
				validation,
				dedupe,
			)
		},
	)

	err := container.Invoke(
		func(
			store storage.Storage,
			hashStore *dedupe.FileHashStore,
			processor *processors.DeduplicationProcessor,
			validation *processors.ValidationProcessor,
			p *pipeline.Pipeline,
		) {
			cfg, err := LoadConfig("./configs/sources.json")
			if err != nil {
				panic(err)
			}

			for _, src := range cfg.Sources {
				f := NewSourceFetcher(src.Name, src.URL)
				if f == nil {
					continue
				}

				source := NewSourceCollector(f, src)
				ctx := context.Background()
				docs, err := source.Collect(ctx)

				if err != nil {
					panic(err)
				}

				for doc := range docs {
					_ = p.Process(doc)
					fmt.Println(doc.Title)
				}
			}
		},
	)

	if err != nil {
		panic(err)
	}
}