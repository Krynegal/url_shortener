package storage

import (
	"fmt"
	"github.com/Krynegal/url_shortener.git/internal/configs"
)

type Storager interface {
	Shorten(string) (int, error)
	Unshorten(string) (string, error)
}

func NewStorage(cfg *configs.Config) (Storager, error) {
	if cfg.FileStorage != "" {
		fs, err := NewFileStorage(cfg.FileStorage)
		if err != nil {
			return nil, err
		}
		if err = fs.ReadURLsFromFile(); err != nil {
			return nil, err
		}
		fmt.Printf("storage: %v", fs.memStorage.store)
		return fs, nil
	}
	return NewMemStorage(), nil
}
