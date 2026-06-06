package storage

import (
	"bufio"
	"os"
)

type FileHashStore struct {
	hashes map[string]struct{}
	path   string
}

func NewFileHashStore(path string) (*FileHashStore, error) {
	store := &FileHashStore{
		hashes: make(map[string]struct{}),
		path:   path,
	}

	file, err := os.Open(path)
	if err == nil {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			store.hashes[scanner.Text()] = struct{}{}
		}
		_ = file.Close()
	}

	return store, nil
}

func (s *FileHashStore) Exists(hash string) bool {
	_, ok := s.hashes[hash]
	return ok
}

func (s *FileHashStore) Add(hash string) error {
	if s.Exists(hash) {
		return nil
	}

	file, err := os.OpenFile(s.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.WriteString(hash + "\n")
	if err != nil {
		return err
	}

	s.hashes[hash] = struct{}{}
	return nil
}
