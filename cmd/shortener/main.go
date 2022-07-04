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

	w.WriteHeader(307)
	w.Header().Set("Location", r.URL.Path)
	log.Println(r.URL.Path)
}

func main() {
	http.HandleFunc("/", ShortURL)
	http.HandleFunc("/get/", GetID)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
