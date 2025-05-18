package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/prakashjegan/app/review-processor/config"
	"github.com/prakashjegan/app/review-processor/review-services/database"
	"github.com/prakashjegan/app/review-processor/review-services/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Run database migrations
	if err := db.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Initialize S3 client
	s3Client, err := utils.NewS3Client(cfg, db)
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start processing files
	go func() {
		for {
			select {
			case <-sigChan:
				log.Println("Shutting down...")
				return
			default:
				// Get unprocessed files
				files, err := s3Client.GetUnprocessedFiles(nil)
				if err != nil {
					log.Printf("Error getting unprocessed files: %v", err)
					continue
				}

				// Process each file
				for _, file := range files {
					// Download file
					data, err := s3Client.DownloadFile(nil, file.S3Key)
					if err != nil {
						log.Printf("Error downloading file %s: %v", file.S3Key, err)
						s3Client.UpdateFileStatus(file.ID, "FAILED", err.Error())
						continue
					}

					// TODO: Add your file processing logic here
					log.Printf("Processing file %s with size %d bytes", file.S3Key, len(data))

					// Update file status
					if err := s3Client.UpdateFileStatus(file.ID, "COMPLETED", ""); err != nil {
						log.Printf("Error updating file status: %v", err)
					}
				}
			}
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
}
