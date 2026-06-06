package pipeline

import (
	"DataCollector/internal/models"
)

type Processor interface {
	Processor(
		example models.Document) (models.Document, bool)
}
