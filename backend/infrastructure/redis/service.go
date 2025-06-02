package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type RedisService interface {
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Delete(ctx context.Context, key string) error
	Ping(ctx context.Context) error
}

type redisService struct {
}

func NewRedisService() RedisService {
	return &redisService{}
}

func (s *redisService) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value for key %s: %w", key, err)
	}

	if err := GetClient().Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set value for key %s: %w", key, err)
	}

	return nil
}

func (s *redisService) Get(ctx context.Context, key string, dest any) error {
	data, err := GetClient().Get(ctx, key).Bytes()
	if err != nil {
		return fmt.Errorf("failed to get value for key %s: %w", key, err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("failed to unmarshal value for key %s: %w", key, err)
	}

	return nil
}

func (s *redisService) Delete(ctx context.Context, key string) error {
	if err := GetClient().Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete key %s: %w", key, err)
	}
	return nil
}

func (s *redisService) Ping(ctx context.Context) error {
	return Ping(ctx)
}
