package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Krynegal/url_shortener.git/internal"
	"github.com/Krynegal/url_shortener.git/internal/configs"
	"github.com/Krynegal/url_shortener.git/internal/handlers/middleware"
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
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

type BatchRequest struct {
	CorrID      string `json:"correlation_id"`
	OriginalURL string `json:"original_url"`
}

type BatchResponse struct {
	CorrID   string `json:"correlation_id"`
	ShortURL string `json:"short_url"`
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
	h.Mux.HandleFunc("/api/shorten/batch", h.Batch).Methods(http.MethodPost)
	h.Mux.HandleFunc("/ping", h.Ping(h.Storage)).Methods(http.MethodGet)
	h.Mux.HandleFunc("/{id}", h.GetID).Methods(http.MethodGet)
	return h
}

func (h *Handler) Ping(st storage.Storager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		db, ok := st.(*storage.DB)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		ctx := r.Context()
		ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func (h *Handler) Batch(w http.ResponseWriter, r *http.Request) {
	session := r.Context().Value(internal.UserIDSessionKey).(internal.Session)

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "can't read body", http.StatusBadRequest)
		return
	}

	var batch = make([]BatchRequest, 0)
	err = json.Unmarshal(body, &batch)
	if err != nil {
		fmt.Printf("Not Decoded: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var resp = make([]BatchResponse, 0)
	for _, v := range batch {
		s, err := h.Storage.Shorten(session.UserID, v.OriginalURL)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		out := fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(s))
		resp = append(resp, BatchResponse{v.CorrID, out})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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

	var status = http.StatusCreated
	newID, err := h.Storage.Shorten(session.UserID, req.URL)
	if err != nil {
		if errors.Is(err, storage.ErrKeyExists) {
			status = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	resp := ResponseAPI{Result: fmt.Sprintf("%s/%s", h.Config.BaseURL, strconv.Itoa(newID))}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
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
	var status = http.StatusCreated
	newID, err := h.Storage.Shorten(session.UserID, url)
	if err != nil {
		if errors.Is(err, storage.ErrKeyExists) {
			status = http.StatusConflict
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	w.WriteHeader(status)
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
