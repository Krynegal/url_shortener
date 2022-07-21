package storage

import (
	"strconv"
	"sync"
)

type Storager interface {
	Shorten(string) int
	Unshorten(string) (string, bool)
}

type Storage struct {
	counter int
	mu      sync.Mutex
	store   map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		counter: 0,
		store:   make(map[string]string),
	}
}

func (s *Storage) Shorten(u string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counter++
	s.store[strconv.Itoa(s.counter)] = u
	return s.counter
}

func (s *Storage) Unshorten(id string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	url, ok := s.store[id]
	return url, ok
}
