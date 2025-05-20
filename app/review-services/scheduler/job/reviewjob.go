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
	aws "github.com/prakashjegan/review-processor/app/review-services/aws"
	"github.com/prakashjegan/review-processor/app/review-services/database/dao"
	"github.com/prakashjegan/review-processor/app/review-services/database/model"
	"github.com/prakashjegan/review-processor/app/review-services/models"
	"github.com/prakashjegan/review-processor/app/review-services/redisstream"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
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
	reviewRawDao := dao.GetReviewRawDao()
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
		totalCount := 0
		totalSuccessFulCount := 0
		if file.Rows != nil {
			totalCount += len(file.Rows)
			for _, row := range file.Rows {
				if row != "" {
					reviewRaw := &model.ReviewRaw{
						ID:                utils.GetUID(),
						ReviewFileStateId: reviewFileState.Id,
						RawData:           row,
						Status:            string(model.RowProcessingStatePending),
						Message:           "",
					}
					reviewRaw, err = reviewRawDao.CreateReviewRaw(reviewRaw)
					if err != nil {
						log.Errorf("Failed to create review raw: %v", err)
						continue
					}
					reviewRawData, _ := json.Marshal(reviewRaw)
					data := map[string]any{}
					data["raw_data"] = reviewRawData
					message := redisstream.Message{
						ID:        fmt.Sprintf("%s-%d-%d", j.jobType, reviewRaw.ID, utils.GetUID()),
						Data:      data,
						Timestamp: time.Now(),
					}
					messagev, err := redisstream.PublishToRedisStream(ctx, message)
					if err != nil {
						log.Errorf("Failed to Publish To Redis Stream: %v", err)
						reviewRaw.Status = string(model.RowProcessingStateFailed)
						reviewRaw.Message = err.Error()
						reviewRaw, err = reviewRawDao.CreateReviewRaw(reviewRaw)

						continue
					}
					reviewRaw.Status = string(model.RowProcessingStateSuccess)
					reviewRaw.Message = messagev
					reviewRaw, err = reviewRawDao.CreateReviewRaw(reviewRaw)
					if err != nil {
						log.Errorf("Failed to update review raw: %v", err)
						continue
					}
					totalSuccessFulCount++
					log.Printf("Review Raw created successfully: %+v\n", reviewRaw)
					//TODO : publish to RedisStreams
				}
			}
		}
		reviewFileState.State = string(models.FileStatusCompleted)
		reviewFileState.TotalRows = totalCount
		reviewFileState.SuccessfulRows = totalSuccessFulCount
		reviewFileState.FailedRows = totalCount - totalSuccessFulCount
		reviewFileState, err = reviewFileStateDao.CreateReviewFileStates(reviewFileState)
		if err != nil {
			log.Errorf("Failed to update review file state: %v", err)
			continue
		}
		// 2. Extract each row and persist it in db. as review raw.
		//3. publish raw review data to kafka topic.

	}

	// Update job event to COMPLETED
	jobEvent.Status = string(models.JobStatusCompleted)
	jobEvent, err = jobEventDao.CreateOrUpdateJobEvent(jobEvent)
	log.Println("Completed review processing job")
}

// Description implements the quartz.Job interface
func (j *ReviewJob) Description() string {
	return "Review processing job"
}

// Key implements the quartz.Job interface
func (j *ReviewJob) Key() int {
	return j.jobKey
}
