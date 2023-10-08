package handler

import (
	"fmt"
	log "github.com/google/logger"
	"github.com/gorilla/mux"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	"io"
	"net/http"
)

func New(store storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

type Handler struct {
	store storage.Storage
}

func (h *Handler) GET(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL, err := h.store.GetURLByID(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
}

func (h *Handler) POST(w http.ResponseWriter, r *http.Request) {
	//Определение схемы
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	longURL := string(b)

	//Проверка, что переданный URl корректный
	if !common.IsValidURL(longURL) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	id, err := h.store.GenIDByURL(longURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if _, tmpErr := w.Write([]byte(fmt.Sprintf("%s://%s/%s", scheme, r.Host, id))); tmpErr != nil {
		log.Error(tmpErr)
	}
}
