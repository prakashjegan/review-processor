package models

import (
	"time"
)

type FileStatus string

const (
	FileStatusPending    FileStatus = "PENDING"
	FileStatusProcessing FileStatus = "PROCESSING"
	FileStatusCompleted  FileStatus = "COMPLETED"
	FileStatusFailed     FileStatus = "FAILED"
)

type File struct {
	ID          int64      `json:"id"`
	S3Key       string     `json:"s3_key"`
	Checksum    string     `json:"checksum"`
	Status      FileStatus `json:"status"`
	RetryCount  int        `json:"retry_count"`
	LastError   string     `json:"last_error"`
	ProcessedAt *time.Time `json:"processed_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

type FileEvent struct {
	ID        int64     `json:"id"`
	FileID    int64     `json:"file_id"`
	EventType string    `json:"event_type"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}
