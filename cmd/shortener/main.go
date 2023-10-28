package main

import (
	"context"
	"errors"
	"github.com/olkonon/shortener/internal/app/config"
	"github.com/olkonon/shortener/internal/app/handler"
	"github.com/olkonon/shortener/internal/app/router"
	"github.com/olkonon/shortener/internal/app/storage/memory"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	appConfig := config.Parse()
	server := &http.Server{
		Handler: router.New(handler.New(memory.NewInMemory(), appConfig.BaseURL)),
		Addr:    appConfig.ListenAddress,
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Infof("Get OS signal [%s], terminating...", sig.String())
		server.SetKeepAlivesEnabled(false)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Error("HTTP server shutdown error: ", err)
		}
	}()
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("Start HTTP server error: ", err)
	}
}
