package dao

import (
	"time"

	"gorm.io/gorm"

	"github.com/prakashjegan/review-processor/app/review-services/models"
)

type ReviewDAO struct {
	db *gorm.DB
}

func NewReviewDAO(db *gorm.DB) *ReviewDAO {
	return &ReviewDAO{db: db}
}

// JobEvent operations
func (d *ReviewDAO) CreateJobEvent(event *models.JobEvent) error {
	return d.db.Create(event).Error
}

func (d *ReviewDAO) UpdateJobEvent(jobID int64, status models.JobStatus, errorMsg string) error {
	updates := map[string]interface{}{
		"status":     status,
		"error":      errorMsg,
		"updated_at": time.Now(),
	}
	if status == models.JobStatusCompleted || status == models.JobStatusFailed {
		now := time.Now()
		updates["completed_at"] = &now
	}

	return d.db.Model(&models.JobEvent{}).Where("id = ?", jobID).Updates(updates).Error
}

// ThirdPartyConfig operations
func (d *ReviewDAO) GetAllThirdPartyConfigs() ([]models.ThirdPartyConfig, error) {
	var configs []models.ThirdPartyConfig
	err := d.db.Find(&configs).Error
	return configs, err
}

// ReviewFileState operations
func (d *ReviewDAO) GetLastProcessedFile(thirdPartyName string) (*models.ReviewFileState, error) {
	var file models.ReviewFileState
	err := d.db.Where("third_party_name = ?", thirdPartyName).
		Order("created_at DESC").
		First(&file).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &file, err
}

func (d *ReviewDAO) CreateFileState(fileState *models.ReviewFileState) error {
	return d.db.Create(fileState).Error
}

func (d *ReviewDAO) UpdateFileState(fileID int64, state string, errorMsg string) error {
	updates := map[string]interface{}{
		"state":      state,
		"error":      errorMsg,
		"updated_at": time.Now(),
	}
	if state == "COMPLETED" {
		now := time.Now()
		updates["processed_at"] = &now
	}

	return d.db.Model(&models.ReviewFileState{}).Where("id = ?", fileID).Updates(updates).Error
}
