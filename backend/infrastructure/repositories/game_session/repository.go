package game_session

import (
	"backend/domain/game_session"
	"backend/pkg/errors"
	"context"
	"fmt"
	"log"
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
	if err := r.db.Create(entity).Error; err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save session", err)
	}
	redisKey := fmt.Sprintf("session:%s:metadata", session.SessionID)
	ctx := context.Background()

	if err := r.redisService.Set(ctx, redisKey, session.Metadata, 2*time.Hour); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save session metadata", err)
	}

	return nil
}

func (r *repository) FindBySessionID(sessionID string) (*game_session.GameSession, error) {
	var entity GameSessionEntity
	if err := r.db.Where("session_id = ?", sessionID).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New(errors.ErrNotFound, "session not found")
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to find session", err)
	}

	session := ToDomain(&entity)

	if session.Status.IsFinished() {
		return nil, errors.New(errors.ErrNotAvailable, "session is no longer active")
	}

	var metadata game_session.SessionMetadata
	redisKey := fmt.Sprintf("session:%s:metadata", sessionID)
	if err := r.redisService.Get(context.Background(), redisKey, &metadata); err != nil {
		// If Redis data not found, mark session as expired
		session.Status = game_session.StatusExpired
		if err := r.db.Model(&GameSessionEntity{}).Where("session_id = ?", session.SessionID).Update("status", session.Status).Error; err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to update expired session status", err)
		}
		return nil, errors.New(errors.ErrNotAvailable, "session has expired")
	}

	session.Metadata = &metadata
	return session, nil
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

func (r *repository) BeginTransaction(sessionID string) (game_session.GameSessionTx, error) {
	// Begin a database transaction
	tx := r.db.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Lock and fetch the session with status check
	var entity GameSessionEntity
	if err := tx.Set("gorm:for_update", true).
		Where("session_id = ? AND status NOT IN (?)", sessionID, []game_session.GameSessionStatus{game_session.StatusFinished, game_session.StatusExpired}).
		First(&entity).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to find active session: %w", err)
	}

	session := ToDomain(&entity)

	var metadata game_session.SessionMetadata
	redisKey := fmt.Sprintf("session:%s:metadata", sessionID)
	if err := r.redisService.Get(context.Background(), redisKey, &metadata); err != nil {
		session.Status = game_session.StatusExpired
		if err := r.db.Model(&GameSessionEntity{}).Where("session_id = ?", session.SessionID).Update("status", session.Status).Error; err != nil {
			log.Printf("failed to update expired session status: %v", err)
		}
		tx.Rollback()
		return nil, fmt.Errorf("failed to get session metadata: %w", err)
	}

	session.Metadata = &metadata
	return &gameSessionTx{
		tx:           tx,
		redisService: r.redisService,
		session:      session,
	}, nil
}

func (r *repository) UpdateGameCraftingStatus(sessionID string, success bool) error {
	var entity GameSessionEntity
	if err := r.db.Where("session_id = ? AND status = ?", sessionID, game_session.StatusStarting).First(&entity).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errors.New(errors.ErrNotFound, "session not found or not in starting status")
		}
		return errors.Wrap(errors.ErrInternal, "failed to find session", err)
	}

	newStatus := game_session.StatusWeek1
	if !success {
		newStatus = game_session.StatusExpired
	}

	if err := r.db.Model(&entity).Update("status", newStatus).Error; err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to update session status", err)
	}

	return nil
}
