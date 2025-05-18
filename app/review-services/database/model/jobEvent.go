package model

import (
	"time"
)

type JobEvent struct {
	Id         int64          `gorm "primary_key" json:"id,omitempty"`
	JobId      int64          `gorm:"index" json:"jobId,omitempty"`
	JobType    string         `json:"jobType,omitempty"`
	EventTime  time.Time      `gorm:"index json:"eventTime,omitempty"`
	EventName  string         `json:"eventName,omitempty"`
	Properties string         `json:"properties,omitempty"`
	Status     string         `json:"status,omitempty"` //TODO : STARTED : SUCCESS : FAILED`
	Message    string         `json:"message,omitempty"`
	CreatedAt  time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}
