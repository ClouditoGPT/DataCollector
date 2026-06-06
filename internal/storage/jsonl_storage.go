package storage

import (
	"DataCollector/internal/models"
	"encoding/json"
	"os"
	"path/filepath"
)

type JsonlStorage struct {
	BasePath string
}

func NewJsonlStorage(basePath string) *JsonlStorage {
	return &JsonlStorage{
		BasePath: basePath,
	}
}

func (s *JsonlStorage) getFilePath(docType models.DocumentType) string {
	return filepath.Join(s.BasePath, string(docType)+".jsonl")
}

func (s *JsonlStorage) Save(
	doc models.Document,
) error {

	path := s.getFilePath(doc.Type)

	if err := os.MkdirAll(
		filepath.Dir(path),
		0755,
	); err != nil {
		return err
	}

	file, err := os.OpenFile(
		path,
		os.O_CREATE|
			os.O_WRONLY|
			os.O_APPEND,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	data, err := json.Marshal(doc)

	if err != nil {
		return err
	}

	_, err = file.Write(
		append(data, '\n'),
	)

	return err
}
