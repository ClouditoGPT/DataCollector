package main

import (
	"DataCollector/internal/dedupe"
	"DataCollector/internal/pipeline"
	"DataCollector/internal/processors"
	"DataCollector/internal/sources/wikipedia"
	"DataCollector/internal/storage"
	"context"
	"fmt"

	"go.uber.org/dig"
)

func main() {
	container := dig.New()

	//? Storages
	container.Provide(func() storage.Storage {
		return storage.NewJsonlStorage("./data")
	})
	container.Provide(func() (*dedupe.FileHashStore, error) {
		return dedupe.NewFileHashStore(
			"./data/hashes.txt",
		)
	})

	//? Processors
	container.Provide(
		processors.NewDeduplicationProcessor,
	)
	container.Provide(
		processors.NewValidationProcessor,
	)

	//? Pipeline
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

	//? Collectors
	container.Provide(
		wikipedia.New,
	)

	err := container.Invoke(
		func(
			store storage.Storage,
			hashStore *dedupe.FileHashStore,
			processor *processors.DeduplicationProcessor,
			validation *processors.ValidationProcessor,
			p *pipeline.Pipeline,
			source *wikipedia.Collector,
		) {

			ctx := context.Background()
			docs, err := source.Collect(ctx)

			if err != nil {
				panic(err)
			}

			for doc := range docs {
				_ = p.Process(doc)

				fmt.Println(doc.Title)
			}
		},
	)

	if err != nil {
		panic(err)
	}
}
