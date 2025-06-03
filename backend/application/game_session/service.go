package game_session

import (
	"backend/domain/game_session"
	"backend/domain/stock"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"time"
)

type Service interface {
	Create(username string, categories []string) (string, error)
	GetState(sessionID string) (*game_session.GameSession, error)
	GetLeaderboard() ([]game_session.GameSession, error)
}

type service struct {
	repo      game_session.Repository
	stockRepo stock.Repository
}

func NewService(repo game_session.Repository, stockRepo stock.Repository) Service {
	return &service{
		repo:      repo,
		stockRepo: stockRepo,
	}
}

func generateSecureToken() (string, error) {
	bytes := make([]byte, 32) // 32 bytes will give us a 64 character hex string
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *service) GetState(sessionID string) (*game_session.GameSession, error) {
	return s.repo.FindBySessionID(sessionID)
}

func (s *service) GetLeaderboard() ([]game_session.GameSession, error) {
	return s.repo.FindLeaderboardTop10(1, 10)
}

func (s *service) Create(username string, categories []string) (string, error) {
	sessionID, err := generateSecureToken()
	if err != nil {
		return "", err
	}

	// Pick stocks for the session
	stocks, err := s.stockRepo.PickStocksForSession(categories)
	if err != nil {
		return "", fmt.Errorf("failed to pick stocks: %w", err)
	}

	// Log picked stocks for debugging
	log.Printf("Picked stocks for session %s: %+v", sessionID, stocks)

	session := &game_session.GameSession{
		SessionID: sessionID,
		Username:  username,
		Cash:      10000.00,
		Status:    game_session.StatusStarting,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	if err := s.repo.Save(session); err != nil {
		return "", err
	}

	return sessionID, nil
}
