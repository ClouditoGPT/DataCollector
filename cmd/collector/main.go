package main

import (
	"DataCollector/internal/dedupe"
	"DataCollector/internal/pipeline"
	"DataCollector/internal/processors"
	"DataCollector/internal/storage"
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
	
	err := container.Invoke(
		func(
			store storage.Storage,
			hashStore *dedupe.FileHashStore,
			processor *processors.DeduplicationProcessor,
			validation *processors.ValidationProcessor,
			p *pipeline.Pipeline,
		) {

			fmt.Println(store)
			fmt.Println(hashStore)
			fmt.Println(processor)
			fmt.Println(validation)
			fmt.Println("pipeline:", p)
		},
	)

	if err != nil {
		panic(err)
	}
}
