package router

import (
	"github.com/gorilla/mux"
	"github.com/olkonon/shortener/internal/app/handler"
	"net/http"
)

func New(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", h.POST).Methods(http.MethodPost)
	r.HandleFunc("/{id}", h.GET).Methods(http.MethodGet)
	return r
}
