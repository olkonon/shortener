package main

import (
	log "github.com/google/logger"
	"github.com/olkonon/shortener/internal/app/handler"
	"net/http"
	"os"
)

const ListenAddress = "127.0.0.1:8080"

func main() {
	log.Init("main", true, true, os.Stdout)
	server := &http.Server{
		Handler: handler.New(),
		Addr:    ListenAddress,
	}
	log.Error(server.ListenAndServe())
}
