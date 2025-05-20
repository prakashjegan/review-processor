package model

import (
	"time"

	"gorm.io/gorm"
)

type ReviewFileStates struct {
	Id              uint64    `gorm:"primaryKey" json:"id,omitempty"`
	FileName        string    `json:"fileName,omitempty"`
	State           string    `json:"state,omitempty"`
	ProcessedDate   time.Time `json:"processedDate,omitempty"`
	FileCreatedDate time.Time `json:"fileCreatedDate,omitempty"`
	FileId          string    `json:"fileId,omitempty"`
	Message         string    `json:"message,omitempty"`

	TotalRows       int            `json:"totalRows,omitempty"`
	SuccessfulRows  int            `json:"successfulRows,omitempty"`
	FailedRows      int            `json:"failedRows,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CheckSum        string         `gorm:"index" json:"checkSum,omitempty"`
	ThirdPartyName  string         `gorm:"index" json:"thirdPartyName,omitempty"`
	ThirdPartyState string         `gorm:"index" json:"thirdPartyState,omitempty"`
}
