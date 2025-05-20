package model

import (
	"time"

	"gorm.io/gorm"
)

type ThirdPartyEvent struct {
	Id            uint64         `gorm:"primaryKey" json:"id,omitempty"`
	JobId         int64          `gorm:"index" json:"jobId,omitempty"`
	JobType       string         `gorm:"index" json:"jobType,omitempty"`
	EventTime     time.Time      `gorm:"index" json:"eventTime,omitempty"`
	ProcessedDate time.Time      `json:"processedDate,omitempty"`
	EventName     string         `json:"eventName,omitempty"`
	Properties    string         `json:"properties,omitempty"`
	Status        string         `json:"status,omitempty"` //TODO : STARTED : SUCCESS : FAILED`
	Message       string         `json:"message,omitempty"`
	CreatedAt     time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted     bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy     string         `gorm:"-" json:"deletedBy,omitempty"`

	ThirdPartyName             string `json:"thirdPartyName,omitempty"`
	ThirdPartyConfigType       string `json:"thirdPartyConfigType,omitempty"`
	ThirdPartyConnectionConfig string `json:"thirdPartyConnectionConfig,omitempty"` // Json Structure
	ThirdPartyReviewConfig     string `json:"thirdPartyReviewConfig,omitempty"`
}
