package cache

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client: client,
	}
}

func (rc *RedisClient) Get(ctx context.Context, key string) (any, error) {
	res, err := rc.client.ZRevRangeWithScores(ctx, key, 0, 4).Result()
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (rc *RedisClient) Set(ctx context.Context, key string, value any) error {
	_, err := rc.client.ZAdd(ctx, key, redis.Z{
		Member: value,
	}).Result()
	if err != nil {
		return err
	}
	return nil
}
