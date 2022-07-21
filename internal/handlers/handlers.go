package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	Mux     *mux.Router
	Storage storage.Storager
	Config  *configs.Config
}

type RequestAPI struct {
	URL string `json:"url"`
}

type ResponseAPI struct {
	Result string `json:"result"`
}

func NewHandler(storage storage.Storager, config *configs.Config) *Handler {
	h := &Handler{
		Mux:     mux.NewRouter(),
		Storage: storage,
		Config:  config,
	}
	h.Mux.HandleFunc("/", h.ShortURL).Methods(http.MethodPost)
	h.Mux.HandleFunc("/api/shorten", h.Shorten).Methods(http.MethodPost)
	h.Mux.HandleFunc("/{id}", h.GetID).Methods(http.MethodGet)
	return h
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	var req RequestAPI
	buf := bytes.NewBuffer(body)
	if err = json.NewDecoder(buf).Decode(&req); err != nil {
		http.Error(w, "wrong request format", http.StatusBadRequest)
		return
	}
	log.Printf("req: %v", req)
	newID := h.Storage.Shorten(req.URL)

	resp := ResponseAPI{Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(newID))}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "troubles with encode response", http.StatusBadRequest)
		return
	}
}

func (h *Handler) ShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	url := string(b)
	newID := h.Storage.Shorten(url)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(newID))))
}

func (h *Handler) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if v, ok := h.Storage.Unshorten(id); ok {
		http.Redirect(w, r, v, http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, "Unknown ID", http.StatusBadRequest)
}
