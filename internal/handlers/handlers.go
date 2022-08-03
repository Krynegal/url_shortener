package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/Krynegal/url_shortener.git/internal"
	"github.com/Krynegal/url_shortener.git/internal/handlers/middleware"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/gorilla/mux"
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

type URL struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type ResponseURLs struct {
	URLs []URL
}

func NewHandler(storage storage.Storager, config *configs.Config) *Handler {
	h := &Handler{
		Mux:     mux.NewRouter(),
		Storage: storage,
		Config:  config,
	}
	h.Mux.Use(middleware.GzipMiddlware, middleware.AuthMiddleware)
	h.Mux.HandleFunc("/", h.ShortURL).Methods(http.MethodPost)
	h.Mux.HandleFunc("/api/shorten", h.Shorten).Methods(http.MethodPost)
	h.Mux.HandleFunc("/api/user/urls", h.GetUrls).Methods(http.MethodGet)
	h.Mux.HandleFunc("/{id}", h.GetID).Methods(http.MethodGet)
	return h
}

func (h *Handler) GetUrls(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(internal.UserIDSessionKey).(internal.Session)
	res := h.Storage.GetAllURLs(session.UserID)
	if len(res) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	u := []URL{}
	for k, v := range res {
		u = append(u, URL{OriginalURL: v, ShortURL: fmt.Sprintf("%s/%s", h.Config.BaseURL, k)})
	}
	fmt.Printf("u: %v\n", u)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(internal.UserIDSessionKey).(internal.Session)

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
	newID, err := h.Storage.Shorten(session.UserID, req.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := ResponseAPI{Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(newID))}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "troubles with encode response", http.StatusBadRequest)
		return
	}
}

func (h *Handler) ShortURL(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(internal.UserIDSessionKey).(internal.Session)

	defer r.Body.Close()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	url := string(b)
	newID, err := h.Storage.Shorten(session.UserID, url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(newID))))
}

func (h *Handler) GetID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	v, err := h.Storage.Unshorten(id)
	if err == nil {
		http.Redirect(w, r, v, http.StatusTemporaryRedirect)
		return
	}
	http.Error(w, err.Error(), http.StatusBadRequest)
}
