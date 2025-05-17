package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/reugn/go-quartz/quartz"
	"zuzu/config"
	"zuzu/review-service/models"
	"zuzu/review-service/s3"
	"zuzu/review-service/storage"
)

type FileProcessorJob struct {
	s3Service *s3.Service
	storage   storage.Storage
	config    *config.Config
}

func NewFileProcessorJob(s3Service *s3.Service, storage storage.Storage, config *config.Config) *FileProcessorJob {
	return &FileProcessorJob{
		s3Service: s3Service,
		storage:   storage,
		config:    config,
	}
}

func (j *FileProcessorJob) Execute(ctx context.Context) error {
	log.Println("Starting file processing job")

	// List all files from S3
	files, err := j.s3Service.ListFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to list files from S3: %v", err)
	}

	for _, s3Key := range files {
		// Check if file is already processed
		existingFile, err := j.storage.GetFileByS3Key(s3Key)
		if err != nil {
			log.Printf("Error checking existing file %s: %v", s3Key, err)
			continue
		}

		if existingFile != nil && existingFile.Status == models.FileStatusCompleted {
			continue
		}

		// Create new file record if not exists
		if existingFile == nil {
			existingFile = &models.File{
				S3Key:      s3Key,
				Status:     models.FileStatusPending,
				RetryCount: 0,
			}
			if err := j.storage.SaveFile(existingFile); err != nil {
				log.Printf("Error saving file record %s: %v", s3Key, err)
				continue
			}
		}

		// Skip if max retries reached
		if existingFile.RetryCount >= j.config.MaxRetries {
			log.Printf("Max retries reached for file %s", s3Key)
			continue
		}

		// Process file
		if err := j.processFile(ctx, existingFile); err != nil {
			log.Printf("Error processing file %s: %v", s3Key, err)
			j.storage.UpdateFileStatus(existingFile.ID, models.FileStatusFailed, err.Error())
			j.storage.SaveFileEvent(&models.FileEvent{
				FileID:    existingFile.ID,
				EventType: "PROCESSING_FAILED",
				Message:   err.Error(),
			})
			existingFile.RetryCount++
			continue
		}

		// Update status to completed
		j.storage.UpdateFileStatus(existingFile.ID, models.FileStatusCompleted, "")
		j.storage.SaveFileEvent(&models.FileEvent{
			FileID:    existingFile.ID,
			EventType: "PROCESSING_COMPLETED",
			Message:   "File processed successfully",
		})
	}

	return nil
}

func (j *FileProcessorJob) processFile(ctx context.Context, file *models.File) error {
	// Download file from S3
	data, checksum, err := j.s3Service.DownloadFile(ctx, file.S3Key)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	// Update file checksum
	file.Checksum = checksum

	// TODO: Add your file processing logic here
	// For example, parse the file content, transform data, etc.
	log.Printf("Processing file %s with checksum %s", file.S3Key, checksum)

	return nil
}

type Scheduler struct {
	scheduler quartz.Scheduler
	job       *FileProcessorJob
}

func NewScheduler(s3Service *s3.Service, storage storage.Storage, config *config.Config) *Scheduler {
	scheduler := quartz.NewStdScheduler()
	job := NewFileProcessorJob(s3Service, storage, config)
	return &Scheduler{
		scheduler: scheduler,
		job:       job,
	}
}

func (s *Scheduler) Start() error {
	s.scheduler.Start()

	// Create a cron trigger
	trigger, err := quartz.NewCronTrigger("0 0/5 * * * ?") // Run every 5 minutes
	if err != nil {
		return fmt.Errorf("failed to create trigger: %v", err)
	}

	// Schedule the job
	err = s.scheduler.ScheduleJob(s.job, trigger)
	if err != nil {
		return fmt.Errorf("failed to schedule job: %v", err)
	}

	return nil
}

func (s *Scheduler) Stop() {
	s.scheduler.Stop()
} 