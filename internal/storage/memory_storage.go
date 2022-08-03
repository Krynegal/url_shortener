package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type MemStorage struct {
	counter   int
	mu        sync.Mutex
	store     map[string]string
	userToIDs map[string][]int
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		counter:   0,
		store:     make(map[string]string),
		userToIDs: make(map[string][]int),
	}
}

func (s *MemStorage) Shorten(uid string, u string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if u == "" {
		return -1, errors.New("url is empty")
	}
	s.counter++
	s.store[strconv.Itoa(s.counter)] = u
	s.userToIDs[uid] = append(s.userToIDs[uid], s.counter)
	fmt.Printf("s.userToIDs[uid]: %v", s.userToIDs)
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

func (s *MemStorage) GetAllURLs(uid string) map[string]string {
	URLSlice := map[string]string{}
	if _, ok := s.userToIDs[uid]; ok {
		for _, id := range s.userToIDs[uid] {
			if url, ok := s.store[strconv.Itoa(id)]; ok {
				URLSlice[strconv.Itoa(id)] = url
			}
		}
	}
	return URLSlice
}
