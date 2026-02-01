package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-job-queue/internal/config"
	"go-job-queue/internal/server"
)

func main() {
	cfg := config.Load()

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: server.NewRouter(),
	}

	shutdownErr := make(chan error, 1)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		shutdownErr <- httpServer.Shutdown(ctx)
	}()

	log.Printf("go-job-queue api listening on :%s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
