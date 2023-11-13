package router

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/olkonon/shortener/internal/app/handler"
	"net/http"
)

func New(h *handler.Handler) *mux.Router {
	r := mux.NewRouter()
	r.Use(handler.WithLog)
	r.Use(handler.WithGzip)
	r.Use(handlers.CompressHandler)
	r.HandleFunc("/", h.POST).Methods(http.MethodPost)
	r.HandleFunc("/ping", h.Ping).Methods(http.MethodGet)
	r.HandleFunc("/{id}", h.GET).Methods(http.MethodGet)
	r.HandleFunc("/api/shorten", h.PostJSON).Methods(http.MethodPost)
	return r
}
