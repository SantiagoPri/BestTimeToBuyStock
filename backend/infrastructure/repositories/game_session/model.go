package game_session

import (
	"backend/domain/game_session"
	"time"
)

type GameSessionEntity struct {
	SessionID string    `gorm:"column:session_id;primaryKey;type:varchar(64)" json:"session_id"`
	Username  string    `gorm:"column:username;type:varchar(100);not null" json:"username"`
	Cash      float64   `gorm:"column:cash;type:decimal(15,2);default:10000.00" json:"cash"`
	Status    string    `gorm:"column:status;type:varchar(20);default:'starting'" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
}

func (GameSessionEntity) TableName() string {
	return "game_sessions"
}

func ToDomain(e *GameSessionEntity) *game_session.GameSession {
	if e == nil {
		return nil
	}
	return &game_session.GameSession{
		SessionID: e.SessionID,
		Username:  e.Username,
		Cash:      e.Cash,
		Status:    game_session.GameSessionStatus(e.Status),
		CreatedAt: e.CreatedAt,
	}
}

func FromDomain(s *game_session.GameSession) *GameSessionEntity {
	if s == nil {
		return nil
	}
	return &GameSessionEntity{
		SessionID: s.SessionID,
		Username:  s.Username,
		Cash:      s.Cash,
		Status:    s.Status.String(),
		CreatedAt: s.CreatedAt,
	}
}
