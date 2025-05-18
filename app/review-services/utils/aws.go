package utils

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	configr "github.com/prakashjegan/review-processor/app/config"
	"github.com/prakashjegan/review-processor/app/review-services/models"
)

type S3Client struct {
	client     *s3.Client
	bucketName string
}

func NewS3Client(cfg *configr.Configuration) (*S3Client, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.Aws.Region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.Aws.AccessKey,
				SecretAccessKey: cfg.Aws.SecreteAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &S3Client{
		client:     client,
		bucketName: cfg.Aws.DocumentBucketName,
	}, nil
}

func (s *S3Client) GetUnprocessedFiles(ctx context.Context) ([]*models.File, error) {
	// Get the last processed file from the database
	lastProcessedFile, err := s.db.GetLastProcessedFile()
	if err != nil {
		return nil, fmt.Errorf("failed to get last processed file: %v", err)
	}

	// List all files from S3
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	var unprocessedFiles []*models.File
	for _, object := range result.Contents {
		// Skip if this file was already processed
		if lastProcessedFile != nil && *object.Key == lastProcessedFile.S3Key {
			break
		}

		// Check if file exists in database
		existingFile, err := s.db.GetFileByS3Key(*object.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to check existing file: %v", err)
		}

		if existingFile == nil {
			// Create new file record
			file := &models.File{
				S3Key:      *object.Key,
				Status:     models.FileStatusPending,
				RetryCount: 0,
				CreatedAt:  time.Now(),
				UpdatedAt:  time.Now(),
			}
			if err := s.db.SaveFile(file); err != nil {
				return nil, fmt.Errorf("failed to save file record: %v", err)
			}
			unprocessedFiles = append(unprocessedFiles, file)
		} else if existingFile.Status == models.FileStatusPending || existingFile.Status == models.FileStatusFailed {
			unprocessedFiles = append(unprocessedFiles, existingFile)
		}
	}

	return unprocessedFiles, nil
}

func (s *S3Client) DownloadFile(ctx context.Context, key string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object body: %v", err)
	}

	return data, nil
}

func (s *S3Client) UpdateFileStatus(fileID int64, status models.FileStatus, errorMsg string) error {
	return s.db.UpdateFileStatus(fileID, status, errorMsg)
}

func (s *S3Client) SaveFileEvent(event *models.FileEvent) error {
	return s.db.SaveFileEvent(event)
}
