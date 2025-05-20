package redisstream

import (
	"context"
	"fmt"
	"time"

	"github.com/prakashjegan/review-processor/app/config"
	"github.com/redis/go-redis/v9"
)

// Publisher handles publishing messages to Redis Stream
type Publisher struct {
	client    *redis.Client
	streamKey string
}

// NewPublisher creates a new Publisher instance
func NewPublisher(client *redis.Client, streamKey string) *Publisher {
	return &Publisher{
		client:    client,
		streamKey: streamKey,
	}
}

// PublishMessage publishes a message to the Redis Stream
func (p *Publisher) PublishMessage(ctx context.Context, message Message) (string, error) {
	// If no ID is provided, use Redis auto-generated ID
	if message.ID == "" {
		message.ID = "*"
	}

	// Convert message data to Redis stream entries
	args := make([]interface{}, 0, len(message.Data)*2+1)
	args = append(args, p.streamKey, message.ID)

	for key, value := range message.Data {
		args = append(args, key, value)
	}

	// Add timestamp if not present
	if message.Timestamp.IsZero() {
		message.Timestamp = time.Now()
	}
	args = append(args, "timestamp", message.Timestamp.Unix())

	// Publish to Redis Stream
	a := &redis.XAddArgs{
		Stream: p.streamKey,
		ID:     message.ID,
		Values: message.Data,
	}
	cmd := p.client.XAdd(ctx, a)
	return cmd.Result()
}

// Close closes the Redis client connection
func (p *Publisher) Close() error {
	return p.client.Close()
}

func PublishToRedisStream(ctx context.Context, message Message) (string, error) {
	redisClient := GetRedisClient()
	publisher := NewPublisher(redisClient, "reviews")
	return publisher.PublishMessage(ctx, message)
}

func GetRedisClient() *redis.Client {
	config := config.GetConfig()
	redisv := config.Database.REDIS

	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisv.Env.Host, redisv.Env.Port),
		Password: "", // no password set
		DB:       0,  // use default DB
	}
	return redis.NewClient(options)
}
