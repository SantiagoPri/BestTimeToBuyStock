package game_session

import (
	"backend/domain/game_session"
	"backend/domain/gm_session"
	"backend/domain/stock"
	"context"
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
	aiModel   gm_session.AI
}

func NewService(repo game_session.Repository, stockRepo stock.Repository, aiModel gm_session.AI) Service {
	return &service{
		repo:      repo,
		stockRepo: stockRepo,
		aiModel:   aiModel,
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

	gmData, err := s.aiModel.GetGMResponse(context.Background(), categories, stocks)
	if err != nil {
		return "", fmt.Errorf("failed to get GM response: %w", err)
	}

	log.Printf("GM Response for session %s: %+v", sessionID, gmData)

	session := &game_session.GameSession{
		SessionID: sessionID,
		Username:  username,
		Cash:      10000.00,
		Status:    game_session.StatusStarting,
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
		Metadata: &game_session.SessionMetadata{
			Holdings: make(map[string]game_session.HoldingInfo),
		},
	}

	if err := s.repo.Save(session); err != nil {
		return "", err
	}

	return sessionID, nil
}
