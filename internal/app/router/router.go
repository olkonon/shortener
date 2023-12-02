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
	r.Use(h.WithAuth)
	r.Methods(http.MethodPost).Path("/").Handler(h.AnonymousAuthHandler(h.POST))
	r.Methods(http.MethodGet).Path("/ping").HandlerFunc(h.Ping)
	r.Methods(http.MethodGet).Path("/{id}").HandlerFunc(h.GET)
	r.Methods(http.MethodPost).Path("/api/shorten/batch").Handler(h.AnonymousAuthHandler(h.BatchPostJSON))
	r.Methods(http.MethodPost).Path("/api/shorten").Handler(h.AnonymousAuthHandler(h.PostJSON))
	r.Methods(http.MethodGet).Path("/api/user/urls").Handler(h.RequireAuthHandler(h.UserGET))
	r.Methods(http.MethodDelete).Path("/api/user/urls").Handler(h.RequireAuthHandler(h.BatchDeleteJSON))
	return r
}
