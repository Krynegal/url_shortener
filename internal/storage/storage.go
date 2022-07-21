package storage

import (
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"strconv"
)

type Storager interface {
	Shorten(string) (int, error)
	Unshorten(string) (string, error)
}

type Storage struct {
	memStore *MemStorage
	file     *FileStorage
}

func NewStorage(cfg *configs.Config) (*Storage, error) {
	memStorage := NewMemStorage()
	var fs *FileStorage
	if cfg.FileStorage != "" {
		var err error
		fs, err = NewFileStorage(cfg.FileStorage)
		if err != nil {
			return nil, err
		}
		err = fs.ReadURLsFromFile(memStorage)
		if err != nil {
			return nil, err
		}
	}
	return &Storage{
		memStore: memStorage,
		file:     fs,
	}, nil
}

func (s *Storage) Shorten(u string) (int, error) {
	v, err := s.memStore.Shorten(u)
	if err != nil {
		return -1, err
	}
	if s.file != nil {
		err := s.file.WriteURLInFile(strconv.Itoa(v), s.memStore.store[strconv.Itoa(v)])
		if err != nil {
			return -1, err
		}
	}
	return v, nil
}

func (s *Storage) Unshorten(id string) (string, error) {
	url, err := s.memStore.Unshorten(id)
	if err != nil {
		return "", err
	}
	return url, nil
}
