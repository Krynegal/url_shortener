package storage

import (
	"log"
	"strconv"
)

type Storage interface {
	Shorten(string) int
	Unshorten(string) (string, bool)
}

type storageType int

const (
	Memory storageType = iota
)

type storage struct {
	counter int
	store   map[string]string
}

func NewStorage(typ storageType) Storage {
	switch typ {
	case Memory:
		return &storage{0, make(map[string]string)}
	default:
		panic("panic")
	}
}

func (s *storage) Shorten(u string) int {
	s.counter++
	s.store[strconv.Itoa(s.counter)] = u
	log.Printf("Shorten storage: %v", s.store)
	return s.counter
}

func (s *storage) Unshorten(id string) (string, bool) {
	url, ok := s.store[id]
	log.Printf("Unshorten storage: %v", s.store)
	return url, ok
}
