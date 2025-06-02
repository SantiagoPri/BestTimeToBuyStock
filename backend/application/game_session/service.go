package game_session

import (
	"backend/domain/game_session"
	"crypto/rand"
	"encoding/hex"
	"time"
)

type Service interface {
	Create(username string) (string, error)
	GetState(sessionID string) (*game_session.GameSession, error)
	UpdateState(sessionID string, newStatus game_session.GameSessionStatus, newCash float64) error
	GetLeaderboard() ([]game_session.GameSession, error)
}

type service struct {
	repo game_session.Repository
}

func NewService(repo game_session.Repository) Service {
	return &service{repo: repo}
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes will give us a 64 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *service) Create(username string) (string, error) {
	sessionID, err := generateSecureToken()
	if err != nil {
		return "", err
	}

	session := &game_session.GameSession{
		SessionID: sessionID,
		Username:  username,
		Cash:      10000.00,
		Status:    game_session.StatusStarting,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Save(session); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (s *service) GetState(sessionID string) (*game_session.GameSession, error) {
	return s.repo.FindBySessionID(sessionID)
}

func (s *service) UpdateState(sessionID string, newStatus game_session.GameSessionStatus, newCash float64) error {
	session, err := s.repo.FindBySessionID(sessionID)
	if err != nil {
		return err
	}

	session.Status = newStatus
	session.Cash = newCash

	return s.repo.Update(session)
}

func (s *service) GetLeaderboard() ([]game_session.GameSession, error) {
	return s.repo.FindLeaderboardTop10(1, 10)
}
