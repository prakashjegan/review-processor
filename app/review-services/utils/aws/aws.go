package utils

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	configr "github.com/prakashjegan/review-processor/app/config"
	"github.com/prakashjegan/review-processor/app/review-services/database/dao"
	"github.com/prakashjegan/review-processor/app/review-services/models"
	"github.com/prakashjegan/review-processor/app/review-services/utils"
	log "github.com/sirupsen/logrus"
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

func (s *S3Client) GetUnprocessedFiles(ctx context.Context, forTp string) ([]*models.File, error) {
	// Get the last processed file from the database
	reviewFileStatedao := dao.GetReviewFileStatesDao()
	lastProcessedFile, err := reviewFileStatedao.GetLastProcessedFile(forTp)
	// List all files from S3
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket:     aws.String(s.bucketName),
		StartAfter: &lastProcessedFile.FileId,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	var unprocessedFiles []*models.File
	for _, object := range result.Contents {
		//TODO : fetch file object.
		data, rows, err := s.DownloadFile(ctx, *object.Key)
		if err != nil {
			log.Debugf("\nError downloading file: %v\n", err)
			//continue
			file := &models.File{
				ID:         utils.GetUID(),
				Name:       *object.Key,
				S3Key:      *object.Key,
				Checksum:   utils.GenerateChecksum(data),
				Status:     models.FileStatusFailed,
				RetryCount: 0,
				LastError:  err.Error(),
				Data:       data,
				Rows:       rows,
			}
			unprocessedFiles = append(unprocessedFiles, file)
		} else {
			file := &models.File{
				ID:         utils.GetUID(),
				Name:       *object.Key,
				S3Key:      *object.Key,
				Checksum:   utils.GenerateChecksum(data),
				Status:     models.FileStatusPending,
				RetryCount: 0,
				LastError:  "",
				Data:       data,
				Rows:       rows,
			}
			unprocessedFiles = append(unprocessedFiles, file)
		}

	}

	return unprocessedFiles, nil
}

func (s *S3Client) DownloadFile(ctx context.Context, key string) ([]byte, []string, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read object body: %v", err)
	}

	var rows []string
	rows = extractRow(data)

	return data, rows, nil
}

func extractRow(data []byte) []string {
	lines := strings.Split(string(data), "\n")
	return lines
}
