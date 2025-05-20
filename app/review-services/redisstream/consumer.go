package redisstream

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Consumer handles consuming messages from Redis Stream
type Consumer struct {
	client     *redis.Client
	streamInfo StreamInfo
}

// NewConsumer creates a new Consumer instance
func NewConsumer(client *redis.Client, streamInfo StreamInfo) *Consumer {
	return &Consumer{
		client:     client,
		streamInfo: streamInfo,
	}
}

// CreateConsumerGroup creates a consumer group if it doesn't exist
func (c *Consumer) CreateConsumerGroup(ctx context.Context) error {
	// Try to create the consumer group
	err := c.client.XGroupCreate(ctx, c.streamInfo.StreamKey, c.streamInfo.GroupName, "0").Err()
	if err != nil && err.Error() != "BUSYGROUP Consumer Group name already exists" {
		return err
	}
	return nil
}

// ConsumeMessages starts consuming messages from the Redis Stream
func (c *Consumer) ConsumeMessages(ctx context.Context, handler func(Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			// Read messages from the stream
			streams, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    c.streamInfo.GroupName,
				Consumer: c.streamInfo.Consumer,
				Streams:  []string{c.streamInfo.StreamKey, ">"},
				Count:    1,
				Block:    0,
			}).Result()

			if err != nil {
				if err == redis.Nil {
					continue
				}
				return err
			}

			// Process messages
			for _, stream := range streams {
				for _, message := range stream.Messages {
					msg := Message{
						ID:        message.ID,
						Data:      message.Values,
						Timestamp: time.Now(),
					}

					// Call the handler function
					if err := handler(msg); err != nil {
						return err
					}

					// Acknowledge the message
					if err := c.client.XAck(ctx, c.streamInfo.StreamKey, c.streamInfo.GroupName, message.ID).Err(); err != nil {
						return err
					}
				}
			}
		}
	}
}

// Close closes the Redis client connection
func (c *Consumer) Close() error {
	return c.client.Close()
}
