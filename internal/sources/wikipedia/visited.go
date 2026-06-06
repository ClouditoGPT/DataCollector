package wikipedia

type Visited struct {
	set map[string]struct{}
}

func NewVisited() *Visited {
	return &Visited{
		set: make(map[string]struct{}),
	}
}

func (v *Visited) Has(url string) bool {
	_, ok := v.set[url]
	return ok
}

func (v *Visited) Add(url string) {
	v.set[url] = struct{}{}
}
