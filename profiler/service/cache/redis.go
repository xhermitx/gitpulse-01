package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient wraps a Redis client to provide simple Get and Set methods.
type RedisClient struct {
	client *redis.Client
}

// NewRedisClient initializes a new RedisClient instance with the provided Redis client.
func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
	}
}

// Get retrieves the value associated with the given key from Redis.
// If the key does not exist, it returns nil and an error if one occurs.
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	res, err := rc.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key '%s' does not exist", key)
	} else if err != nil {
		return "", fmt.Errorf("failed to get key '%s': %w", key, err)
	}
	return res, nil
}

// Set sets a key-value pair in Redis with an expiration time (TTL).
// The key will expire after the provided duration (t), or never if t is 0.
func (rc *RedisClient) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	_, err := rc.client.Set(ctx, key, value, expiration).Result()
	if err != nil {
		return fmt.Errorf("failed to set key '%s': %w", key, err)
	}
	return nil
}

func (rc *RedisClient) Append(ctx context.Context, key string, value string) error {
	res := rc.client.Append(ctx, key, value)
	if res.Err() != nil {
		return res.Err()
	}
	return nil
}
