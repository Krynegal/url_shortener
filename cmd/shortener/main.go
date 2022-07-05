package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
)

var urls = make(map[string]string)

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
	//log.Println(r.Header.Get("Host"))
	shortUrl := strings.Split(string(b), "/")
	urls[string(b)] = shortUrl[len(shortUrl)-1]
	//log.Println(urls)
	w.WriteHeader(201)
	w.Write(b)
}

func GetID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
		return
	}
	vars := mux.Vars(r)
	//log.Println(vars)
	id := vars["id"]
	if id == "" {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	//log.Println(id)
	for k, v := range urls {
		if id == v {
			w.Header().Set("Location", k)
			w.WriteHeader(307)
		}
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ShortURL)
	r.HandleFunc("/{id}", GetID)
	log.Fatal(http.ListenAndServe("localhost:8080", r))
}
