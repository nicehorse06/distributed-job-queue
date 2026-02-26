package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"distributed-job-queue/internal/computeclient"
	"distributed-job-queue/internal/config"
	computev1 "distributed-job-queue/internal/gen/compute/v1"
	"distributed-job-queue/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cfg := config.Load()

	computeConn, err := grpc.Dial(cfg.ComputeAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to compute service %q: %v", cfg.ComputeAddr, err)
	}
	defer computeConn.Close()

	computeSvc := computeclient.New(computev1.NewComputeServiceClient(computeConn))

	httpServer := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: server.NewRouter(computeSvc),
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

	log.Printf("distributed-job-queue api listening on :%s", cfg.Port)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	if err := <-shutdownErr; err != nil {
		log.Printf("graceful shutdown error: %v", err)
	}
}
