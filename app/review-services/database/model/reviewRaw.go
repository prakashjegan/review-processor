package model

import (
	"time"

	"gorm.io/gorm"
)

type ReviewRaw struct {
	ID                uint64 `gorm:"primaryKey" json:"id,omitempty"`
	ReviewFileStateId uint64 `gorm:"index" json:"reviewFileStateId,omitempty"`

	RawData       string    `json:"rawData,omitempty"`
	Status        string    `json:"status,omitempty"`
	Message       string    `json:"message,omitempty"`
	ProcessedDate time.Time `json:"processedDate,omitempty"`

	CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"createdAt,omitempty"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime;index" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	IsDeleted bool           `gorm:"-" json:"isDeleted,omitempty"`
	DeletedBy string         `gorm:"-" json:"deletedBy,omitempty"`
}
