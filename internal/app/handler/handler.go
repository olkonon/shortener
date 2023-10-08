package handler

import (
	"fmt"
	log "github.com/google/logger"
	"github.com/olkonon/shortener/internal/app/common"
	"github.com/olkonon/shortener/internal/app/storage"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	"io"
	"net/http"
	"path"
)

func New() *Handler {
	return &Handler{
		store: memory.NewInMemory(),
	}
}

type Handler struct {
	store storage.Storage
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Error(err)
		}
	}()

	switch r.Method {
	case http.MethodPost:
		//Проверка URL если не подошло 404
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		//Определение схемы запроса
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

		//Проверка корректности URL
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
			log.Error(err)
		}

	case http.MethodGet:
		if path.Dir(r.URL.Path) != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		id := path.Base(r.URL.Path)
		longURL, err := h.store.GetURLByID(id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		http.Redirect(w, r, longURL, http.StatusTemporaryRedirect)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
