package crawler

import (
	"encoding/json"
	"os"
)

type VisitedStore struct {
	path string
}

func NewVisitedStore(path string) *VisitedStore {
	return &VisitedStore{
		path: path,
	}
}

func (s *VisitedStore) Save(set map[string]struct{}) error {
	list := make([]string, 0, len(set))
	for k := range set {
		list = append(list, k)
	}
	data, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0644)
}

func (s *VisitedStore) Load() (map[string]struct{}, error) {
	data, err := os.ReadFile(s.path)
	if err != nil {
		return map[string]struct{}{}, nil
	}
	var list []string
	_ = json.Unmarshal(data, &list)
	set := make(map[string]struct{})
	for _, v := range list {
		set[v] = struct{}{}
	}
	return set, nil
}