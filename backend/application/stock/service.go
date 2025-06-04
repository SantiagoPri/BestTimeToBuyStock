package stockapp

import (
	"backend/domain/stock"
	"context"
)

type StockService struct {
	repo stock.Repository
}

func NewStockService(repo stock.Repository) *StockService {
	return &StockService{repo: repo}
}

func (s *StockService) FindOne(field string, value any) (*stock.Stock, error) {
	filters := map[string]any{field: value}
	return s.repo.FindBy(filters)
}

func (s *StockService) FindAllStocks(ctx context.Context, page, limit int, filters map[string]string, sortBy, sortOrder string) ([]stock.Stock, int64, error) {
	params := stock.NewQueryParams()
	params.Page = page
	params.PageSize = limit
	params.Filters = filters
	params.SortBy = sortBy
	params.SortOrder = sortOrder
	return s.repo.FindAllStocks(ctx, params)
}
