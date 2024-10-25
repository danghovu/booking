package main

import (
	"log"

	"booking-event/config"
	"booking-event/internal/app/worker"
)

func main() {
	config, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	server := worker.NewServer(*config)
	if err := server.Run(); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
