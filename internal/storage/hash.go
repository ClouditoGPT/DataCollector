package storage

type HashStore interface {
	Exists(hash string) bool
	Add(hash string) error
}
