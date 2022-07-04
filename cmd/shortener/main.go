package main

import (
	"io"
	"log"
	"net/http"
	"strings"
)

func ShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", 400)
		return
	}
	log.Println(string(b))
	w.WriteHeader(201)
	w.Write(b)
}

func GetID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/get/")
	if id == "" {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	log.Println(id)
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	log.Println(r.URL.Path)

	w.Header().Set("Location", "http://localhost:8080/")
	w.WriteHeader(307)
}

func main() {
	http.HandleFunc("/", ShortURL)
	http.HandleFunc("/get/", GetID)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
