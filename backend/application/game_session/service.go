package game_session

import (
	gmsvc "backend/application/gm_session"
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
	Buy(sessionID string, ticker string, quantity int) error
	Sell(sessionID string, ticker string, quantity int) error
	AdvanceWeek(sessionID string) error
	EndSession(sessionID string) error
}

type service struct {
	repo      game_session.Repository
	stockRepo stock.Repository
	aiModel   gm_session.AI
	gmService gmsvc.Service
}

func NewService(
	repo game_session.Repository,
	stockRepo stock.Repository,
	aiModel gm_session.AI,
	gmService gmsvc.Service,
) Service {
	return &service{
		repo:      repo,
		stockRepo: stockRepo,
		aiModel:   aiModel,
		gmService: gmService,
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

	stocks, err := s.stockRepo.PickStocksForSession(categories)
	if err != nil {
		return "", fmt.Errorf("failed to pick stocks: %w", err)
	}

	gmData, err := s.aiModel.GetGMResponse(context.Background(), categories, stocks)
	if err != nil {
		return "", fmt.Errorf("failed to get GM response: %w", err)
	}

	log.Printf("GM Response for session %s: %+v", sessionID, gmData)

	if err := s.gmService.SaveGMWeekData(sessionID, gmData); err != nil {
		return "", fmt.Errorf("failed to save GM week data: %w", err)
	}

	initialCash := 10000.00
	session := &game_session.GameSession{
		SessionID:     sessionID,
		Username:      username,
		Cash:          initialCash,
		HoldingsValue: 0.00,
		TotalBalance:  initialCash,
		Status:        game_session.StatusWeek1,
		CreatedAt:     time.Now().Format(time.RFC3339),
		UpdatedAt:     time.Now().Format(time.RFC3339),
		Metadata: &game_session.SessionMetadata{
			Holdings: make(map[string]game_session.HoldingInfo),
		},
	}

	if err := s.repo.Save(session); err != nil {
		return "", err
	}

	return sessionID, nil
}

func getCurrentWeek(status game_session.GameSessionStatus) (int, error) {
	switch status {
	case game_session.StatusWeek1:
		return 1, nil
	case game_session.StatusWeek2:
		return 2, nil
	case game_session.StatusWeek3:
		return 3, nil
	case game_session.StatusWeek4:
		return 4, nil
	case game_session.StatusWeek5:
		return 5, nil
	default:
		return 0, fmt.Errorf("invalid game status for trading: %s", status)
	}
}

func (s *service) Buy(sessionID string, ticker string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	tx, err := s.repo.BeginTransaction(sessionID)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	session := tx.GetSession()

	currentWeek, err := getCurrentWeek(session.Status)
	if err != nil {
		return err
	}

	gmData, err := s.gmService.GetWeekData(sessionID, currentWeek)
	if err != nil {
		return fmt.Errorf("failed to get GM week data: %w", err)
	}

	var stockPrice float64
	found := false
	for _, stock := range gmData.Stocks {
		if stock.Ticker == ticker {
			stockPrice = stock.Price
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("stock %s not found in current week data", ticker)
	}

	totalCost := stockPrice * float64(quantity)

	if totalCost > session.Cash {
		return fmt.Errorf("insufficient funds: need %.2f, have %.2f", totalCost, session.Cash)
	}

	if session.Metadata == nil {
		session.Metadata = &game_session.SessionMetadata{
			Holdings: make(map[string]game_session.HoldingInfo),
		}
	}

	holding, exists := session.Metadata.Holdings[ticker]
	if exists {
		holding.Quantity += quantity
		holding.TotalSpent += totalCost
	} else {
		holding = game_session.HoldingInfo{
			Quantity:   quantity,
			TotalSpent: totalCost,
		}
	}
	session.Metadata.Holdings[ticker] = holding

	session.Cash -= totalCost

	holdingsValue := 0.0
	for ticker, holding := range session.Metadata.Holdings {
		for _, stock := range gmData.Stocks {
			if stock.Ticker == ticker {
				holdingsValue += float64(holding.Quantity) * stock.Price
				break
			}
		}
	}

	session.HoldingsValue = holdingsValue
	session.TotalBalance = session.Cash + session.HoldingsValue
	session.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := tx.Update(session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *service) Sell(sessionID string, ticker string, quantity int) error {
	if quantity <= 0 {
		return fmt.Errorf("quantity must be positive")
	}

	tx, err := s.repo.BeginTransaction(sessionID)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	session := tx.GetSession()

	if session.Metadata == nil || session.Metadata.Holdings == nil {
		return fmt.Errorf("no holdings found")
	}

	holding, exists := session.Metadata.Holdings[ticker]
	if !exists {
		return fmt.Errorf("no holdings found for stock %s", ticker)
	}

	if holding.Quantity < quantity {
		return fmt.Errorf("insufficient stocks: have %d, want to sell %d", holding.Quantity, quantity)
	}

	currentWeek, err := getCurrentWeek(session.Status)
	if err != nil {
		return err
	}

	gmData, err := s.gmService.GetWeekData(sessionID, currentWeek)
	if err != nil {
		return fmt.Errorf("failed to get GM week data: %w", err)
	}

	var stockPrice float64
	found := false
	for _, stock := range gmData.Stocks {
		if stock.Ticker == ticker {
			stockPrice = stock.Price
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("stock %s not found in current week data", ticker)
	}

	saleProceeds := stockPrice * float64(quantity)

	holding.Quantity -= quantity
	spentPerShare := holding.TotalSpent / float64(holding.Quantity+quantity)
	holding.TotalSpent -= spentPerShare * float64(quantity)
	session.Metadata.Holdings[ticker] = holding
	session.Cash += saleProceeds

	holdingsValue := 0.0
	for ticker, holding := range session.Metadata.Holdings {
		for _, stock := range gmData.Stocks {
			if stock.Ticker == ticker {
				holdingsValue += float64(holding.Quantity) * stock.Price
				break
			}
		}
	}

	session.HoldingsValue = holdingsValue
	session.TotalBalance = session.Cash + session.HoldingsValue
	session.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := tx.Update(session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *service) AdvanceWeek(sessionID string) error {
	tx, err := s.repo.BeginTransaction(sessionID)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	session := tx.GetSession()

	currentWeek, err := getCurrentWeek(session.Status)
	if err != nil {
		return err
	}

	if currentWeek >= 5 {
		return fmt.Errorf("cannot advance beyond week 5")
	}

	var nextStatus game_session.GameSessionStatus
	switch session.Status {
	case game_session.StatusWeek1:
		nextStatus = game_session.StatusWeek2
	case game_session.StatusWeek2:
		nextStatus = game_session.StatusWeek3
	case game_session.StatusWeek3:
		nextStatus = game_session.StatusWeek4
	case game_session.StatusWeek4:
		nextStatus = game_session.StatusWeek5
	default:
		return fmt.Errorf("invalid game status for advancing week: %s", session.Status)
	}

	nextWeek := currentWeek + 1
	gmData, err := s.gmService.GetWeekData(sessionID, nextWeek)
	if err != nil {
		return fmt.Errorf("failed to get GM week data: %w", err)
	}

	holdingsValue := 0.0
	for ticker, holding := range session.Metadata.Holdings {
		for _, stock := range gmData.Stocks {
			if stock.Ticker == ticker {
				holdingsValue += float64(holding.Quantity) * stock.Price
				break
			}
		}
	}

	session.Status = nextStatus
	session.HoldingsValue = holdingsValue
	session.TotalBalance = session.Cash + session.HoldingsValue
	session.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := tx.Update(session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *service) EndSession(sessionID string) error {
	tx, err := s.repo.BeginTransaction(sessionID)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	session := tx.GetSession()

	currentWeek, err := getCurrentWeek(session.Status)
	if err != nil {
		return err
	}

	if currentWeek != 5 {
		return fmt.Errorf("can only end session in week 5, current week: %d", currentWeek)
	}

	gmData, err := s.gmService.GetWeekData(sessionID, currentWeek)
	if err != nil {
		return fmt.Errorf("failed to get GM week data: %w", err)
	}

	for ticker, holding := range session.Metadata.Holdings {
		var stockPrice float64
		for _, stock := range gmData.Stocks {
			if stock.Ticker == ticker {
				stockPrice = stock.Price
				break
			}
		}

		saleProceeds := stockPrice * float64(holding.Quantity)
		session.Cash += saleProceeds
	}

	session.Metadata.Holdings = make(map[string]game_session.HoldingInfo)
	session.HoldingsValue = 0
	session.TotalBalance = session.Cash
	session.Status = game_session.StatusFinished
	session.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := tx.Update(session); err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
