package stockapp

import "backend/domain/stock"

type StockService struct {
	repo stock.Repository
}

func NewStockService(repo stock.Repository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) FindPaginated(page, limit int) ([]stock.Stock, int64, error) {
	return s.repo.FindPaginated(page, limit)
}

func (s *StockService) FindOne(field string, value any) (*stock.Stock, error) {
	filters := map[string]any{field: value}
	return s.repo.FindBy(filters)
}
