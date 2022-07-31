package storage

import (
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
		return fs, nil
	}
	return NewMemStorage(), nil
}
