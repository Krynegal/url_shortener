package main

import (
	"log"
	"net/http"
	"strings"
)

func HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
	q := r.URL.Query().Get("query")
	log.Printf("search string: %s", q)
	if q == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}
	// обработка входящего url

	//if r.Method != http.MethodPost {
	//	http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
	//	return
	//}
}

func GetId(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("get method"))
	id := strings.TrimPrefix(r.URL.Path, "/get/")
	if id == "" {
		http.Error(w, "wrong id", http.StatusBadRequest)
		return
	}
	log.Println(id)
	//if r.Method != http.MethodGet {
	//	http.Error(w, "Only GET requests are allowed", http.StatusMethodNotAllowed)
	//	return
	//}
}

func main() {
	//http.HandleFunc("/", HelloWorld)
	//http.HandleFunc("/get/", GetId)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
