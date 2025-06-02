package game_session

import (
	"backend/domain/game_session"
	"context"
	"fmt"
	"time"

	"backend/infrastructure/redis"

	"gorm.io/gorm"
)

type repository struct {
	db           *gorm.DB
	redisService redis.RedisService
}

func NewRepository(db *gorm.DB, redisService redis.RedisService) game_session.Repository {
	return &repository{
		db:           db,
		redisService: redisService,
	}
}

func (r *repository) Save(session *game_session.GameSession) error {
	entity := FromDomain(session)
	return r.db.Create(entity).Error
}

func (r *repository) FindBySessionID(sessionID string) (*game_session.GameSession, error) {
	var entity GameSessionEntity
	if err := r.db.Where("session_id = ?", sessionID).First(&entity).Error; err != nil {
		return nil, err
	}

	session := ToDomain(&entity)

	// Only fetch metadata for active sessions
	if session.Status != game_session.StatusFinished && session.Status != game_session.StatusExpired {
		var metadata game_session.SessionMetadata
		redisKey := fmt.Sprintf("session:%s:metadata", sessionID)

		if err := r.redisService.Get(context.Background(), redisKey, &metadata); err == nil {
			session.Metadata = &metadata
		}
	}

	return session, nil
}

func (r *repository) Update(session *game_session.GameSession) error {
	// Update database fields
	entity := FromDomain(session)
	if err := r.db.Save(entity).Error; err != nil {
		return err
	}

	if session.Metadata != nil && session.Status != game_session.StatusFinished && session.Status != game_session.StatusExpired {
		redisKey := fmt.Sprintf("session:%s:metadata", session.SessionID)
		ctx := context.Background()

		if err := r.redisService.Set(ctx, redisKey, session.Metadata, 2*time.Hour); err != nil {
			return fmt.Errorf("failed to update session metadata: %w", err)
		}
	}

	return nil
}

func (r *repository) FindLeaderboardTop10(page, pageSize int) ([]game_session.GameSession, error) {
	var entities []GameSessionEntity
	offset := (page - 1) * pageSize

	if err := r.db.Where("status = ?", game_session.StatusFinished).
		Order("cash DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	sessions := make([]game_session.GameSession, len(entities))
	for i, entity := range entities {
		sessions[i] = *ToDomain(&entity)
	}
	return sessions, nil
}
