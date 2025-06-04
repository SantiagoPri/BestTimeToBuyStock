package gm_session

import (
	"backend/domain/gm_session"
	"fmt"
)

type Service interface {
	SaveGMWeekData(sessionID string, gmData map[string]*gm_session.GMWeekData) error
	GetWeekData(sessionID string, week int) (*gm_session.GMWeekData, error)
}

type service struct {
	repo gm_session.Repository
}

func NewService(repo gm_session.Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) SaveGMWeekData(sessionID string, gmData map[string]*gm_session.GMWeekData) error {
	// Save data for each week (1-5)
	for i := 1; i <= 5; i++ {
		weekKey := fmt.Sprintf("week%d", i)
		weekData, exists := gmData[weekKey]
		if !exists {
			return fmt.Errorf("missing data for %s", weekKey)
		}

		if err := s.repo.SaveWeekData(sessionID, i, weekData); err != nil {
			return fmt.Errorf("failed to save data for %s: %w", weekKey, err)
		}
	}

	return nil
}

func (s *service) GetWeekData(sessionID string, week int) (*gm_session.GMWeekData, error) {
	if week < 1 || week > 5 {
		return nil, fmt.Errorf("invalid week number: %d, must be between 1 and 5", week)
	}
	return s.repo.GetWeekData(sessionID, week)
}
