package handlers

import (
	"github.com/Krynegal/url_shortener.git/internal/storage"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

type Handler struct {
	Mux     *mux.Router
	Storage storage.IStorage
}

func NewHandler(storage storage.IStorage) *Handler {
	h := &Handler{
		Mux:     mux.NewRouter(),
		Storage: storage,
	}
	h.Mux.Handle("/", h.ShortURL()).Methods(http.MethodPost)
	h.Mux.Handle("/{id}", h.GetID()).Methods(http.MethodGet)
	return h
}

func (h *Handler) ShortURL() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		url := string(b)
		newID := h.Storage.Shorten(url)
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("http://localhost:8080/" + strconv.Itoa(newID)))
	}
}

func (h *Handler) GetID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		if v, ok := h.Storage.Unshorten(id); ok {
			http.Redirect(w, r, v, http.StatusTemporaryRedirect)
			return
		}
		http.Error(w, "Unknown ID", http.StatusBadRequest)
	}
}
