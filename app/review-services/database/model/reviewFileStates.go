package model

import (
	"time"

	"gorm.io/gorm"
)

type ReviewFileStates struct {
	Id              int64          `gorm:"primaryKey" json:"id,omitempty"`
	FileName        string         `json:"fileName,omitempty"`
	State           string         `json:"state,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CheckSum        string         `gorm:"index" json:"checkSum,omitempty"`
	ThirdPartyName  string         `gorm:"index" json:"thirdPartyName,omitempty"`
	ThirdPartyState string         `gorm:"index" json:"thirdPartyState,omitempty"`
}
