package gm_session

import (
	"backend/domain/gm_session"
	"backend/pkg/errors"
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
			return errors.New(errors.ErrInvalidInput, "missing data for "+weekKey)
		}

		if err := s.repo.SaveWeekData(sessionID, i, weekData); err != nil {
			return errors.Wrap(errors.ErrInternal, "failed to save data for "+weekKey, err)
		}
	}

	return nil
}

func (s *service) GetWeekData(sessionID string, week int) (*gm_session.GMWeekData, error) {
	if week < 1 || week > 5 {
		return nil, errors.New(errors.ErrInvalidInput, "invalid week number: must be between 1 and 5")
	}
	return s.repo.GetWeekData(sessionID, week)
}
