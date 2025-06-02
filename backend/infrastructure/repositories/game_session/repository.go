package game_session

import (
	"backend/domain/game_session"

	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) game_session.Repository {
	return &repository{db: db}
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
	return ToDomain(&entity), nil
}

func (r *repository) Update(session *game_session.GameSession) error {
	entity := FromDomain(session)
	return r.db.Save(entity).Error
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
