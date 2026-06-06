package pipeline

import (
	"DataCollector/internal/models"
	"DataCollector/internal/storage"
)

type Pipeline struct {
	processors []Processor
	storage    storage.Storage
}

func NewPipeline(store storage.Storage, processor ...Processor) *Pipeline {
	return &Pipeline{
		processors: processor,
		storage:    store,
	}
}

func (p *Pipeline) Process(doc models.Document) error {
	current := doc
	for _, processor := range p.processors {
		next, ok := processor.Process(current)
		if !ok {
			return nil
		}
		current = next
	}

	return p.storage.Save(current)
}
