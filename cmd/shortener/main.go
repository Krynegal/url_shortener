package main

import (
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/handlers"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := configs.Get()
	s, err := storage.NewStorage(cfg)
	if err != nil {
		panic("can't create storage")
	}
	r := handlers.NewHandler(s, cfg).Mux
	log.Printf("server run on address: %s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
