package main

import (
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/handlers"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := configs.GetConfigs()
	s := storage.NewStorage()
	r := handlers.NewHandler(s, cfg).Mux
	log.Printf("%s", cfg.ServerAddress)
	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}
