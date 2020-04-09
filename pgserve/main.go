package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"allaboutapps.at/aw/go-mranftl-sample/pgserve/api"
	"allaboutapps.at/aw/go-mranftl-sample/pgserve/router"
	"allaboutapps.at/aw/go-mranftl-sample/pgtestpool"
)

func main() {
	manager := pgtestpool.DefaultManagerFromEnv()
	if err := manager.Initialize(context.Background()); err != nil {
		log.Fatalf("Failed to initialize testpool manager: %v", err)
	}

	server := &api.Server{M: manager}
	router := router.Init(server)

	go func() {
		if err := router.Start(":8080"); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := router.Shutdown(ctx); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to gracefully shut down HTTP server: %v", err)
	}
}
