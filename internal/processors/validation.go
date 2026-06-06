package processors

import (
	"fmt"

	"DataCollector/internal/models"
)

type ValidationProcessor struct {
	minLength int
}

func NewValidationProcessor() *ValidationProcessor {
	return &ValidationProcessor{
		minLength: 100,
	}
}

func (p *ValidationProcessor) Process(
	doc models.Document,
) (models.Document, bool) {

	text := fmt.Sprint(doc.Content)

	if len(text) < p.minLength {
		return doc, false
	}

	return doc, true
}
