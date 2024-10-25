package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"booking-event/config"
	"booking-event/internal/app/auth"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	server := auth.NewServer(*config)
	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to run server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-quit
	log.Println("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), config.GracefulShutdown)
	defer cancel()

	done := make(chan struct{})
	go func() {
		defer close(done)
		server.Shutdown(context.Background())
	}()

	select {
	case <-done:
		log.Println("Server shutdown gracefully")
	case <-ctx.Done():
		log.Fatalf("Server shutdown timed out")
	}
}
