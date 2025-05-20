package redisstream

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisConfig holds the configuration for Redis connection
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// NewRedisClient creates a new Redis client with the given configuration
func NewRedisClient(config RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}

// Message represents a message to be published to Redis Stream
type Message struct {
	ID        string
	Data      map[string]interface{}
	Timestamp time.Time
}

// StreamInfo holds information about the Redis Stream
type StreamInfo struct {
	StreamKey string
	GroupName string
	Consumer  string
}
