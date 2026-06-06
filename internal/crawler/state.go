package crawler

import (
	"sync"
	"time"
)

type State struct {
	mu       sync.RWMutex
	source   string
	running  bool
	visited  int
	queue    int
	errors   int
	started  time.Time
}

var states = make(map[string]*State)
var statesMu sync.Mutex

func GetState(source string) *State {
	statesMu.Lock()
	defer statesMu.Unlock()
	if s, ok := states[source]; ok {
		return s
	}
	s := &State{source: source}
	states[source] = s
	return s
}

func (s *State) SetRunning(r bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = r
	if r {
		s.started = time.Now()
	}
}

func (s *State) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *State) SetVisited(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.visited = n
}

func (s *State) GetVisited() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.visited
}

func (s *State) SetQueue(n int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue = n
}

func (s *State) GetQueue() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue
}

func (s *State) IncErrors() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.errors++
}

func (s *State) GetErrors() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.errors
}

func (s *State) GetUptime() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if !s.running {
		return 0
	}
	return time.Since(s.started)
}