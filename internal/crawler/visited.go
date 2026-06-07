package crawler

import "sync"

type Visited struct {
	mu  sync.Mutex
	set map[string]struct{}
}

func NewVisited() *Visited {
	return &Visited{
		set: make(map[string]struct{}),
	}
}

func (v *Visited) Has(url string) bool {
	v.mu.Lock()
	defer v.mu.Unlock()
	_, ok := v.set[url]
	return ok
}

func (v *Visited) Add(url string) {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.set[url] = struct{}{}
}

func (v *Visited) Snapshot() map[string]struct{} {
	v.mu.Lock()
	defer v.mu.Unlock()
	copy := make(map[string]struct{}, len(v.set))
	for k := range v.set {
		copy[k] = struct{}{}
	}
	return copy
}