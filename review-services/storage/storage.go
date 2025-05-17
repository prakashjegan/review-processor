package storage

import (
	"sync"
	"time"

	"zuzu/review-service/models"
)

type Storage interface {
	SaveFile(file *models.File) error
	GetFileByS3Key(s3Key string) (*models.File, error)
	UpdateFileStatus(fileID int64, status models.FileStatus, errorMsg string) error
	SaveFileEvent(event *models.FileEvent) error
	GetUnprocessedFiles() ([]*models.File, error)
}

type InMemoryStorage struct {
	mu     sync.RWMutex
	files  map[int64]*models.File
	events []*models.FileEvent
	nextID int64
}

func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{
		files:  make(map[int64]*models.File),
		events: make([]*models.FileEvent, 0),
		nextID: 1,
	}
}

func (s *InMemoryStorage) SaveFile(file *models.File) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file.ID = s.nextID
	s.nextID++
	file.CreatedAt = time.Now()
	file.UpdatedAt = time.Now()
	s.files[file.ID] = file
	return nil
}

func (s *InMemoryStorage) GetFileByS3Key(s3Key string) (*models.File, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, file := range s.files {
		if file.S3Key == s3Key {
			return file, nil
		}
	}
	return nil, nil
}

func (s *InMemoryStorage) UpdateFileStatus(fileID int64, status models.FileStatus, errorMsg string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, exists := s.files[fileID]
	if !exists {
		return nil
	}

	file.Status = status
	file.LastError = errorMsg
	file.UpdatedAt = time.Now()
	if status == models.FileStatusCompleted {
		now := time.Now()
		file.ProcessedAt = &now
	}
	return nil
}

func (s *InMemoryStorage) SaveFileEvent(event *models.FileEvent) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event.ID = s.nextID
	s.nextID++
	event.CreatedAt = time.Now()
	s.events = append(s.events, event)
	return nil
}

func (s *InMemoryStorage) GetUnprocessedFiles() ([]*models.File, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var unprocessed []*models.File
	for _, file := range s.files {
		if file.Status == models.FileStatusPending || file.Status == models.FileStatusFailed {
			unprocessed = append(unprocessed, file)
		}
	}
	return unprocessed, nil
} 