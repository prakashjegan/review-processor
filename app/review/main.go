package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"zuzu/config"
	"zuzu/review-service/s3"
	"zuzu/review-service/scheduler"
	"zuzu/review-service/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize S3 service
	s3Service, err := s3.NewService(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize S3 service: %v", err)
	}

	// Initialize storage
	storage := storage.NewInMemoryStorage()

	// Initialize and start scheduler
	scheduler := scheduler.NewScheduler(s3Service, storage, cfg)
	if err := scheduler.Start(); err != nil {
		log.Fatalf("Failed to start scheduler: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	log.Println("Shutting down...")
	scheduler.Stop()
} 