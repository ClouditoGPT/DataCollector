package storage

import "DataCollector/internal/models"

type Storage interface {
	Save(doc models.Document) error
}
