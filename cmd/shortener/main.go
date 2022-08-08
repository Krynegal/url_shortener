package main

import (
	"flag"
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/handlers"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := configs.NewConfig()

	flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "Server address")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "Base URL")
	flag.StringVar(&cfg.FileStorage, "f", cfg.FileStorage, "File Storage Path")
	flag.StringVar(&cfg.DB, "d", cfg.DB, "Data base path string")
	flag.Parse()

	s, err := storage.NewStorage(cfg)
	if err != nil {
		panic("can't create storage")
	}
	r := handlers.NewHandler(s, cfg).Mux
	log.Printf("server run on address: %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
