package scheduler

import (
	"context"
	"log"

	"github.com/prakashjegan/review-processor/app/config"
)

// ReviewJob handles scheduled review processing tasks
type ReviewJob struct {
	config *config.Configuration
	jobKey int
}

func NewReviewJob(config *config.Configuration) *ReviewJob {
	return &ReviewJob{
		config: config,
		jobKey: 1, // You might want to generate a unique key
	}
}

// Execute implements the quartz.Job interface
func (j *ReviewJob) Execute(ctx context.Context) {
	log.Println("Starting file processing job")
	// db := gdatabase.GetDB()
	// // List all files from S3
	// client := utils.NewS3Client(j.config, db)
	// if err != nil {
	// 	log.Printf("failed to list files from S3: %v", err)
	// 	return
	// }

	// for _, s3Key := range files {
	// 	// Check if file is already processed
	// 	existingFile, err := j.storage.GetFileByS3Key(s3Key)
	// 	if err != nil {
	// 		log.Printf("Error checking existing file %s: %v", s3Key, err)
	// 		continue
	// 	}

	// 	if existingFile != nil && existingFile.Status == models.FileStatusCompleted {
	// 		continue
	// 	}

	// 	// Create new file record if not exists
	// 	if existingFile == nil {
	// 		existingFile = &models.File{
	// 			S3Key:      s3Key,
	// 			Status:     models.FileStatusPending,
	// 			RetryCount: 0,
	// 		}
	// 		if err := j.storage.SaveFile(existingFile); err != nil {
	// 			log.Printf("Error saving file record %s: %v", s3Key, err)
	// 			continue
	// 		}
	// 	}

	// 	// Skip if max retries reached
	// 	if existingFile.RetryCount >= j.config.MaxRetries {
	// 		log.Printf("Max retries reached for file %s", s3Key)
	// 		continue
	// 	}

	// 	// Process file
	// 	if err := j.processFile(ctx, existingFile); err != nil {
	// 		log.Printf("Error processing file %s: %v", s3Key, err)
	// 		j.storage.UpdateFileStatus(existingFile.ID, models.FileStatusFailed, err.Error())
	// 		j.storage.SaveFileEvent(&models.FileEvent{
	// 			FileID:    existingFile.ID,
	// 			EventType: "PROCESSING_FAILED",
	// 			Message:   err.Error(),
	// 		})
	// 		existingFile.RetryCount++
	// 		continue
	// 	}

	// 	// Update status to completed
	// 	j.storage.UpdateFileStatus(existingFile.ID, models.FileStatusCompleted, "")
	// 	j.storage.SaveFileEvent(&models.FileEvent{
	// 		FileID:    existingFile.ID,
	// 		EventType: "PROCESSING_COMPLETED",
	// 		Message:   "File processed successfully",
	// 	})
	// }

	log.Println("Completed file processing job")
}

// Description implements the quartz.Job interface
func (j *ReviewJob) Description() string {
	return "Review processing job"
}

// Key implements the quartz.Job interface
func (j *ReviewJob) Key() int {
	return j.jobKey
}

// func (j *FileProcessorJob) processFile(ctx context.Context, file *models.File) error {
// 	// Download file from S3
// 	data, checksum, err := j.s3Service.DownloadFile(ctx, file.S3Key)
// 	if err != nil {
// 		return fmt.Errorf("failed to download file: %v", err)
// 	}

// 	// Update file checksum
// 	file.Checksum = checksum

// 	// TODO: Add your file processing logic here
// 	// For example, parse the file content, transform data, etc.
// 	log.Printf("Processing file %s with checksum %s", file.S3Key, checksum)

// 	return nil
// }
