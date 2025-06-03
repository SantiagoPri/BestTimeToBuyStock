package game_session

import (
	"backend/domain/game_session"
	"backend/infrastructure/redis"
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type gameSessionTx struct {
	tx           *gorm.DB
	redisService redis.RedisService
	session      *game_session.GameSession
}

func (tx *gameSessionTx) GetSession() *game_session.GameSession {
	return tx.session
}

func (tx *gameSessionTx) Update(session *game_session.GameSession) error {
	// Update database fields
	entity := FromDomain(session)
	if err := tx.tx.Save(entity).Error; err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if session.Metadata != nil && !session.Status.IsFinished() {
		redisKey := fmt.Sprintf("session:%s:metadata", session.SessionID)
		ctx := context.Background()

		if err := tx.redisService.Set(ctx, redisKey, session.Metadata, 2*time.Hour); err != nil {
			return fmt.Errorf("failed to update session metadata: %w", err)
		}
	}

	tx.session = session
	return nil
}

func (tx *gameSessionTx) Commit() error {
	return tx.tx.Commit().Error
}

func (tx *gameSessionTx) Rollback() error {
	return tx.tx.Rollback().Error
}
