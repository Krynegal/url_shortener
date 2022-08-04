package configs

import (
	"github.com/caarlos0/env/v6"
	"log"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage   string `env:"FILE_STORAGE_PATH"`
	DB            string `env:"DB" envDefault:"postgres://localhost:5432/urlshorten?sslmode=disable"`
}

func NewConfig() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		log.Fatal(err)
	}
	log.Printf("configs: %v", *cfg)
	return cfg
}
