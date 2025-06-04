package redis

import (
	"backend/pkg/errors"
	"context"
	"encoding/json"
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
		return errors.Wrap(errors.ErrInternal, "failed to marshal value", err)
	}

	if err := GetClient().Set(ctx, key, data, ttl).Err(); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to set value in Redis", err)
	}

	return nil
}

func (s *redisService) Get(ctx context.Context, key string, dest any) error {
	data, err := GetClient().Get(ctx, key).Bytes()
	if err != nil {
		return errors.Wrap(errors.ErrNotFound, "failed to get value from Redis", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to unmarshal value", err)
	}

	return nil
}

func (s *redisService) Delete(ctx context.Context, key string) error {
	if err := GetClient().Del(ctx, key).Err(); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete key from Redis", err)
	}
	return nil
}

func (s *redisService) Ping(ctx context.Context) error {
	if err := Ping(ctx); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to ping Redis", err)
	}
	return nil
}
