package main

import (
	"github.com/Krynegal/url_shortener.git/internal/handlers"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"log"
	"net/http"
)

func main() {
	s := storage.NewStorage()
	r := handlers.NewHandler(s).Mux
	log.Fatal(http.ListenAndServe(":8080", r))
}
