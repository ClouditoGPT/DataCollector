package crawler

import "sync"

type Queue struct {
	mu    sync.Mutex
	items []string
}

func NewQueue(seed []string) *Queue {
	return &Queue{
		items: seed,
	}
}

func (q *Queue) Push(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, url)
}

func (q *Queue) Pop() (string, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return "", false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

func (q *Queue) Snapshot() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := make([]string, len(q.items))
	copy(out, q.items)
	return out
}