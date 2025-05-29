package stock

import (
	"backend/domain/stock"
	"backend/infrastructure/repositories"

	"gorm.io/gorm"
)

type StockRepository struct {
	repo repositories.Repository
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{
		repo: repositories.NewBaseRepository(db, &stock.Stock{}),
	}
}

func (r *StockRepository) Save(s stock.Stock) error {
	return r.repo.Save(&s)
}

func (r *StockRepository) FindAll() ([]stock.Stock, error) {
	var stocks []stock.Stock
	err := r.repo.FindAll(&stocks)
	if err != nil {
		return nil, err
	}
	return stocks, nil
}

func (r *StockRepository) DeleteByTicker(ticker string) error {
	return r.repo.DeleteByField("ticker", ticker, &stock.Stock{})
}

func (r *StockRepository) FindBy(filters map[string]any) (*stock.Stock, error) {
	var s stock.Stock
	err := r.repo.FindOneBy(filters, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *StockRepository) FindPaginated(page int, limit int) ([]stock.Stock, int64, error) {
	var stocks []stock.Stock
	total, err := r.repo.FindPaginated(&stocks, page, limit)
	if err != nil {
		return nil, 0, err
	}
	return stocks, total, nil
}
