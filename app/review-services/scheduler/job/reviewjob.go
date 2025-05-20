package scheduler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/prakashjegan/review-processor/app/config"
	"github.com/prakashjegan/review-processor/app/database"
	gdatabase "github.com/prakashjegan/review-processor/app/database"
	"github.com/prakashjegan/review-processor/app/review-services/database/dao"
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/models"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
	aws "github.com/prakashjegan/review-processor/app/review-services/utils/aws"
)

// ReviewJob handles scheduled review processing tasks
type ReviewJob struct {
	config  *config.Configuration
	jobKey  int
	jobType string
	db      *gorm.DB
}

func NewReviewJob(jobType string, config *config.Configuration) *ReviewJob {
	return &ReviewJob{
		config:  config,
		jobKey:  1,
		jobType: jobType,
		db:      database.GetDB(),
	}
}

// Execute implements the quartz.Job interface
func (j *ReviewJob) Execute(ctx context.Context) {
	log.Println("Starting review processing job")
	db := gdatabase.GetDB()
	// Create job event with STARTED status
	jobId := utils.GetUID()
	jobEvent := &model.JobEvent{
		// Id            uint64         `gorm:"primaryKey" json:"id,omitempty"`
		// JobId         int64          `gorm:"index" json:"jobId,omitempty"`
		// JobType       string         `gorm:"index" json:"jobType,omitempty"`
		// EventTime     time.Time      `gorm:"index" json:"eventTime,omitempty"`
		// ProcessedDate time.Time      `json:"processedDate,omitempty"`
		// EventName     string         `json:"eventName,omitempty"`
		// Properties    string         `json:"properties,omitempty"`
		// Status        string         `json:"status,omitempty"` //TODO : STARTED : SUCCESS : FAILED`
		// Message       string         `json:"message,omitempty"`
		// CreatedAt     time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
		// UpdatedAt     time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
		// DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
		// IsDeleted     bool           `gorm:"-" json:"isDeleted,omitempty"`
		// DeletedBy     string         `gorm:"-" json:"deletedBy,omitempty"`

		Id:        jobId,
		JobId:     jobId,
		JobType:   "REVIEW_PROCESSING_JOB",
		EventTime: time.Now(),
		EventName: "REVIEW_PROCESSING_JOB",
		Status:    string(models.JobStatusStarted),
	}
	jobEventDao := dao.GetJobEventDao()
	reviewFileStateDao := dao.GetReviewFileStatesDao()
	jobEvent, err := jobEventDao.CreateOrUpdateJobEvent(jobEvent)
	if err != nil {
		log.Errorf("Failed to create job event: %v", err)
		return
	}

	// Fetch all third party configs
	var thirdPartyConfigs []model.ThirdPartyConfig
	err = db.Where("third_party_name = ? ", j.jobType).Find(&thirdPartyConfigs).Error
	if err != nil {
		log.Errorf("Failed to fetch third party configs: %v", err)
		// jobEvent.Status = string(models.JobStatusFailed)
		// jobEvent, err = jobEventDao.CreateOrUpdateJobEvent(jobEvent)
		// return
	}
	// TODO :
	// 1. Fetch Unprocessed Files from aws.go file function.
	s3Client, err := aws.NewS3Client(
		j.config,
	)
	if err != nil {
		log.Errorf("Failed to create S3 client: %v", err)
		jobEvent.Status = string(models.JobStatusFailed)
		jobEvent.Message = err.Error()
		jobEvent, err = jobEventDao.CreateOrUpdateJobEvent(jobEvent)
		log.Println("Completed review processing job")
		return
	}
	files, err := s3Client.GetUnprocessedFiles(context.Background(), j.jobType)
	if err != nil {
		log.Errorf("Failed to Fetch Unprocessed Files from S3 client: %v", err)
		jobEvent.Status = string(models.JobStatusFailed)
		jobEvent.Message = err.Error()
		jobEvent, err = jobEventDao.CreateOrUpdateJobEvent(jobEvent)
		log.Println("Completed review processing job")
		return
	}

	// 1. Persist data in reviewFileStates.go
	for _, file := range files {
		reviewFileState := &model.ReviewFileStates{
			Id:              utils.GetUID(),
			FileName:        file.S3Key,
			FileCreatedDate: file.CreatedAt,
			FileId:          file.S3Key,
			CheckSum:        file.Checksum,
			State:           string(models.FileStatusProcessing),
			ThirdPartyState: string(models.FileStatusProcessing),
			ThirdPartyName:  j.jobType,
		}
		reviewFileState, err = reviewFileStateDao.CreateReviewFileStates(reviewFileState)
		if err != nil {
			log.Errorf("Failed to create review file states: %v", err)
			continue
		}
		log.Printf("Review File State created successfully: %+v\n", reviewFileState)
		if file.Rows != nil {

			for _, row := range file.Rows {

			}
		}
		// 2. Extract each row and persist it in db. as review raw.
		//3. publish raw review data to kafka topic.

	}

	// Update job event to COMPLETED
	jobEvent.Status = string(models.JobStatusCompleted)
	jobEvent, err = jobEventDao.CreateOrUpdateJobEvent(jobEvent)
	log.Println("Completed review processing job")
}

func (j *ReviewJob) processThirdParty(ctx context.Context, config models.ThirdPartyConfig) error {
	// Get last processed file for this third party
	var lastProcessedFile models.ReviewFileState
	if err := j.db.Where("third_party_name = ?", config.Name).
		Order("created_at DESC").
		First(&lastProcessedFile).Error; err != nil && err != gorm.ErrRecordNotFound {
		return fmt.Errorf("failed to get last processed file: %v", err)
	}

	// Initialize S3 client for this third party
	s3Client, err := utils.NewS3Client(&config)
	if err != nil {
		return fmt.Errorf("failed to create S3 client: %v", err)
	}

	// List files from S3
	files, err := s3Client.ListFiles(ctx)
	if err != nil {
		return fmt.Errorf("failed to list files: %v", err)
	}

	// Process each file
	for _, s3Key := range files {
		// Skip if file was already processed
		if lastProcessedFile.S3Key == s3Key {
			continue
		}

		// Create file state record
		fileState := &models.ReviewFileState{
			ThirdPartyName: config.Name,
			S3Key:          s3Key,
			State:          "DOWNLOADED",
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		if err := j.db.Create(fileState).Error; err != nil {
			log.Errorf("Failed to create file state for %s: %v", s3Key, err)
			continue
		}

		// Download and process file
		if err := j.processFile(ctx, s3Client, fileState); err != nil {
			log.Errorf("Failed to process file %s: %v", s3Key, err)
			j.updateFileState(fileState.ID, "FAILED", err.Error())
			continue
		}

		// Update file state to completed
		j.updateFileState(fileState.ID, "COMPLETED", "")
	}

	return nil
}

func (j *ReviewJob) processFile(ctx context.Context, s3Client *utils.S3Client, fileState *models.ReviewFileState) error {
	// Download file
	data, err := s3Client.DownloadFile(ctx, fileState.S3Key)
	if err != nil {
		return fmt.Errorf("failed to download file: %v", err)
	}

	// Process each line as a JSON review
	var reviews []models.ReviewBO
	if err := json.Unmarshal(data, &reviews); err != nil {
		return fmt.Errorf("failed to parse reviews: %v", err)
	}

	// Process each review
	for _, review := range reviews {
		if err := j.processReview(review); err != nil {
			log.Errorf("Failed to process review %s: %v", review.ReviewID, err)
			continue
		}
	}

	return nil
}

func (j *ReviewJob) processReview(review models.ReviewBO) error {
	// Start a transaction
	tx := j.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// TODO: Implement review processing logic
	// 1. Check if product exists, if not create
	// 2. Check if customer exists, if not create
	// 3. Create review record
	// 4. Update any related statistics

	return tx.Commit().Error
}

func (j *ReviewJob) updateJobEvent(jobID int64, status models.JobStatus, errorMsg string) {
	updates := map[string]interface{}{
		"status":     status,
		"error":      errorMsg,
		"updated_at": time.Now(),
	}
	if status == models.JobStatusCompleted || status == models.JobStatusFailed {
		now := time.Now()
		updates["completed_at"] = &now
	}

	if err := j.db.Model(&models.JobEvent{}).Where("id = ?", jobID).Updates(updates).Error; err != nil {
		log.Errorf("Failed to update job event: %v", err)
	}
}

func (j *ReviewJob) updateFileState(fileID int64, state string, errorMsg string) {
	updates := map[string]interface{}{
		"state":      state,
		"error":      errorMsg,
		"updated_at": time.Now(),
	}
	if state == "COMPLETED" {
		now := time.Now()
		updates["processed_at"] = &now
	}

	if err := j.db.Model(&models.ReviewFileState{}).Where("id = ?", fileID).Updates(updates).Error; err != nil {
		log.Errorf("Failed to update file state: %v", err)
	}
}

// Description implements the quartz.Job interface
func (j *ReviewJob) Description() string {
	return "Review processing job"
}

// Key implements the quartz.Job interface
func (j *ReviewJob) Key() int {
	return j.jobKey
}
