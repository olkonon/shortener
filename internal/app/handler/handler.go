package handler

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
)

func New(store storage.Storage, baseURL string) *Handler {
	return &Handler{
		store:   store,
		baseURL: baseURL,
	}
}

type Handler struct {
	store   storage.Storage
	baseURL string
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
	if _, tmpErr := w.Write([]byte(fmt.Sprintf("%s/%s", h.baseURL, id))); tmpErr != nil {
		log.Error(tmpErr)
	}
}
