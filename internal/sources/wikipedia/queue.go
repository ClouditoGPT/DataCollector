package wikipedia

type Queue struct {
	items []string
}

func NewQueue(seed []string) *Queue {
	return &Queue{
		items: seed,
	}
}

func (q *Queue) Push(url string) {
	q.items = append(q.items, url)
}

func (q *Queue) Pop() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}
