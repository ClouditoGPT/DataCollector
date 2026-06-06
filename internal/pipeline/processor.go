package pipeline

import (
	"DataCollector/internal/models"
)

type Processor interface {
	Process(
		doc models.Document) (models.Document, bool)
}
