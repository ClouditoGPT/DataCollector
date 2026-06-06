package crawler

import (
	"encoding/json"
	"os"
)

type QueueStore struct {
	path string
}

func NewQueueStore(path string) *QueueStore {
	return &QueueStore{
		path: path,
	}
}

func (s *QueueStore) Save(items []string) error {
	data, err := json.Marshal(&items)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *QueueStore) Load() ([]string, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, err
	}
	var items []string
	err = json.Unmarshal(data, &items)
	return items, err
}