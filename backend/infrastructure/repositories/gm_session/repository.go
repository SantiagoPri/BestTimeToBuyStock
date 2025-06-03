package gm_session

import (
	"context"
	"fmt"
	"time"

	"backend/domain/gm_session"
	"backend/infrastructure/redis"
)

type repository struct {
	redisService redis.RedisService
}

func NewRepository(redisService redis.RedisService) gm_session.Repository {
	return &repository{
		redisService: redisService,
	}
}

func (r *repository) SaveWeekData(sessionID string, week int, data *gm_session.GMWeekData) error {
	key := fmt.Sprintf("gm:session:%s:week:%d", sessionID, week)
	return r.redisService.Set(context.Background(), key, data, 2*time.Hour)
}

func (r *repository) GetWeekData(sessionID string, week int) (*gm_session.GMWeekData, error) {
	key := fmt.Sprintf("gm:session:%s:week:%d", sessionID, week)
	var data gm_session.GMWeekData
	err := r.redisService.Get(context.Background(), key, &data)
	if err != nil {
		return nil, fmt.Errorf("failed to get week data: %w", err)
	}
	return &data, nil
}

func (r *repository) ClearSessionData(sessionID string) error {
	ctx := context.Background()
	for week := 1; week <= 5; week++ {
		key := fmt.Sprintf("gm:session:%s:week:%d", sessionID, week)
		if err := r.redisService.Delete(ctx, key); err != nil {
			continue
		}
	}
	return nil
}
