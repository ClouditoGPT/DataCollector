package processors

import (
	"DataCollector/internal/dedupe"
	"DataCollector/internal/models"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

type DeduplicationProcessor struct {
	store dedupe.HashStore
}

func NewDeduplicationProcessor(
	store *dedupe.FileHashStore,
) *DeduplicationProcessor {

	return &DeduplicationProcessor{
		store: store,
	}
}

func getText(
	doc models.Document,
) string {
	return fmt.Sprint(
		doc.Content,
	)
}

func (p *DeduplicationProcessor) Process(
	doc models.Document,
) (models.Document, bool) {

	text := getText(doc)

	sum := sha256.Sum256(
		[]byte(text),
	)

	hash := hex.EncodeToString(
		sum[:],
	)

	if p.store.Exists(hash) {
		return doc, false
	}

	_ = p.store.Add(hash)

	return doc, true
}
