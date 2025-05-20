package model

import (
	"time"

	"gorm.io/gorm"
)

type ThirdPartyConfig struct {
	ID                         uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ThirdPartyName             string `json:"thirdPartyName,omitempty"`
	ThirdPartyConfigType       string `json:"thirdPartyConfigType,omitempty"`
	ThirdPartyConnectionConfig string `json:"thirdPartyConnectionConfig,omitempty"` // Json Structure
	ThirdPartyReviewConfig     string `json:"thirdPartyReviewConfig,omitempty"`     // Json Structure

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}

// Config :
