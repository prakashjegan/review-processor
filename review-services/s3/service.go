package s3

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"zuzu/config"
)

type Service struct {
	client     *s3.Client
	bucketName string
}

func NewService(cfg *config.Config) (*Service, error) {
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWSRegion),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{
				AccessKeyID:     cfg.AWSAccessKeyID,
				SecretAccessKey: cfg.AWSSecretAccessKey,
			}, nil
		})),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(awsCfg)
	return &Service{
		client:     client,
		bucketName: cfg.S3BucketName,
	}, nil
}

func (s *Service) ListFiles(ctx context.Context) ([]string, error) {
	result, err := s.client.ListObjectsV2(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects: %v", err)
	}

	var files []string
	for _, object := range result.Contents {
		files = append(files, *object.Key)
	}
	return files, nil
}

func (s *Service) DownloadFile(ctx context.Context, key string) ([]byte, string, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, "", fmt.Errorf("failed to get object: %v", err)
	}
	defer result.Body.Close()

	data, err := io.ReadAll(result.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to read object body: %v", err)
	}

	// Calculate MD5 checksum
	hash := md5.Sum(data)
	checksum := hex.EncodeToString(hash[:])

	return data, checksum, nil
} 