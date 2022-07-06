package main

import (
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

var urls = make(map[string]string)
var id int

func ShortURL(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	url := string(b)
	//key := genShortenURL(url)
	key := strconv.Itoa(id)
	urls[key] = url
	id++
	log.Println(urls)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("http://localhost:8080/" + key))
}

//func genShortenURL(s string) string {
//	u := strconv.Itoa(id)
//	id++
//	return fmt.Sprintf("http://%s/%s", "localhost:8080", u)
//}

func GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	log.Println(id)
	if v, ok := urls[id]; ok {
		http.Redirect(w, r, v, http.StatusTemporaryRedirect)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusBadRequest)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", ShortURL).Methods(http.MethodPost)
	r.HandleFunc("/{id}", GetID).Methods(http.MethodGet)
	r.HandleFunc("/", DefaultHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
