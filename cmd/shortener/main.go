package main

import (
	log "github.com/google/logger"
	"github.com/olkonon/shortener/internal/app/config"
	"github.com/olkonon/shortener/internal/app/handler"
	"github.com/olkonon/shortener/internal/app/router"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	"net/http"
	"os"
)

func main() {
	log.Init("main", true, true, os.Stdout)
	appConfig := config.Parse()
	server := &http.Server{
		Handler: router.New(handler.New(memory.NewInMemory(), appConfig.BaseURL)),
		Addr:    appConfig.ListenAddress,
	}
	log.Error(server.ListenAndServe())
}
