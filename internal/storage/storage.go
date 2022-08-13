package storage

import (
	"errors"
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/storage/postgres"
)

var ErrKeyExists = errors.New("url is already shorten")

type Storager interface {
	Shorten(string, string) (int, error)
	Unshorten(string) (string, error)
	GetAllURLs(string) map[string]string
}

func NewStorage(cfg *configs.Config) (Storager, error) {
	if cfg.DB != "" {
		postgresDB, err := postgres.NewDatabaseStorage(cfg.DB)
		if err != nil {
			return nil, err
		}

		db := NewDB(postgresDB)
		return db, nil
	}
	if cfg.FileStorage != "" {
		fs, err := NewFileStorage(cfg.FileStorage)
		if err != nil {
			return nil, err
		}
		return fs, nil
	}
	return NewMemStorage(), nil
}
