package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	// S3 Configuration
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	S3BucketName       string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Scheduler Configuration
	SchedulerInterval string
	MaxRetries       int
	RetryDelay       time.Duration

	// Application Configuration
	LogLevel     string
	Environment  string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	return &Config{
		AWSRegion:          os.Getenv("AWS_REGION"),
		AWSAccessKeyID:     os.Getenv("AWS_ACCESS_KEY_ID"),
		AWSSecretAccessKey: os.Getenv("AWS_SECRET_ACCESS_KEY"),
		S3BucketName:       os.Getenv("S3_BUCKET_NAME"),

		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),

		SchedulerInterval: os.Getenv("SCHEDULER_INTERVAL_1"),
		MaxRetries:       3, // Default value
		RetryDelay:       time.Duration(300) * time.Second,

		LogLevel:    os.Getenv("LOG_LEVEL"),
		Environment: os.Getenv("ENVIRONMENT"),
	}, nil
} 