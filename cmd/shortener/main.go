package main

import (
	"context"
	log "github.com/google/logger"
	"github.com/olkonon/shortener/internal/app/config"
	"github.com/olkonon/shortener/internal/app/handler"
	"github.com/olkonon/shortener/internal/app/router"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Init("main", true, true, os.Stdout)
	appConfig := config.Parse()
	server := &http.Server{
		Handler: router.New(handler.New(memory.NewInMemory(), appConfig.BaseURL)),
		Addr:    appConfig.ListenAddress,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Infof("Get OS signal %v terminating...", sig)
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("HTTP server shutdown error: ", err)
		}
	}()
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Error("Start HTTP server error: ", err)
	}
}
