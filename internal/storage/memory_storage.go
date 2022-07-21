package storage

import (
	"errors"
	"strconv"
	"sync"
)

type MemStorage struct {
	counter int
	mu      sync.Mutex
	store   map[string]string
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter: 0,
		store:   make(map[string]string),
	}
}

func (s *MemStorage) Shorten(u string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if u == "" {
		return -1, errors.New("url is empty")
	}
	s.counter++
	s.store[strconv.Itoa(s.counter)] = u
	return s.counter, nil
}

func (s *MemStorage) Unshorten(id string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if url, ok := s.store[id]; ok {
		return url, nil
	}
	return "", errors.New("id doesn't exist")
}
