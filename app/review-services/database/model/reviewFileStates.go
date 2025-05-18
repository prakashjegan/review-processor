package model

import (
	"time"

	"gorm.io/gorm"
)

type ReviewFileStates struct {
	Id              int64          `gorm:"primary_key" json:"id,omitempty"`
	FileName        string         `json:"fileName,omitempty"`
	State           string         `json:"state,omitempty"`
	CreatedAt       time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	CheckSum        string         `json:"checkSum,omitempty"`
	ThirdPartyName  string         `json:"thirdPartyName,omitempty"`
	ThirdPartyState string         `json:"thirdPartyState,omitempty"`
}
