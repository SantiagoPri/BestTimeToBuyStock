package stock

import "context"

type Repository interface {
	FindAllStocks(ctx context.Context, params QueryParams) ([]Stock, int64, error)
	FindBy(filters map[string]any) (*Stock, error)
	PickStocksForSession(categories []string) ([]Stock, error)
}
