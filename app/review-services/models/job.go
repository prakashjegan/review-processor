package models

type JobStatus string

const (
	JobStatusStarted   JobStatus = "STARTED"
	JobStatusCompleted JobStatus = "COMPLETED"
	JobStatusFailed    JobStatus = "FAILED"
)

// type JobEvent struct {
// 	ID          int64     `gorm:"primaryKey"`
// 	JobName     string    `gorm:"not null"`
// 	Status      JobStatus `gorm:"not null"`
// 	Error       string
// 	StartedAt   time.Time `gorm:"not null"`
// 	CompletedAt *time.Time
// 	CreatedAt   time.Time `gorm:"not null"`
// 	UpdatedAt   time.Time `gorm:"not null"`
// }

// type ThirdPartyConfig struct {
// 	ID              int64     `gorm:"primaryKey"`
// 	Name            string    `gorm:"not null;unique"`
// 	BucketName      string    `gorm:"not null"`
// 	Region          string    `gorm:"not null"`
// 	AccessKey       string    `gorm:"not null"`
// 	SecretAccessKey string    `gorm:"not null"`
// 	CreatedAt       time.Time `gorm:"not null"`
// 	UpdatedAt       time.Time `gorm:"not null"`
// }

// type ReviewFileState struct {
// 	ID             int64  `gorm:"primaryKey"`
// 	ThirdPartyName string `gorm:"not null"`
// 	S3Key          string `gorm:"not null"`
// 	State          string `gorm:"not null"`
// 	Error          string
// 	ProcessedAt    *time.Time
// 	CreatedAt      time.Time `gorm:"not null"`
// 	UpdatedAt      time.Time `gorm:"not null"`
// }

// type ReviewBO struct {
// 	ReviewID         string    `json:"review_id"`
// 	ProductID        string    `json:"product_id"`
// 	CustomerID       string    `json:"customer_id"`
// 	Rating           int       `json:"rating"`
// 	Title            string    `json:"title"`
// 	Content          string    `json:"content"`
// 	ReviewDate       time.Time `json:"review_date"`
// 	VerifiedPurchase bool      `json:"verified_purchase"`
// 	HelpfulVotes     int       `json:"helpful_votes"`
// 	TotalVotes       int       `json:"total_votes"`
// 	Language         string    `json:"language"`
// 	Source           string    `json:"source"`
// }
