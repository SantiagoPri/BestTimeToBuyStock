package stock

import (
	"backend/domain/stock"
	"backend/infrastructure/repositories"

	"gorm.io/gorm"
)

// Compile-time assertion to ensure StockRepository implements stock.Repository
var _ stock.Repository = (*StockRepository)(nil)

type StockRepository struct {
	repo repositories.Repository
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{
		repo: repositories.NewBaseRepository(db, &StockEntity{}),
	}
}

func (r *StockRepository) Save(s *stock.Stock) error {
	entity := FromDomain(s)
	return r.repo.Save(entity)
}

func (r *StockRepository) FindAll() ([]stock.Stock, error) {
	var entities []StockEntity
	err := r.repo.FindAll(&entities)
	if err != nil {
		return nil, err
	}

	stocks := make([]stock.Stock, len(entities))
	for i, entity := range entities {
		domainStock := ToDomain(&entity)
		stocks[i] = *domainStock
	}
	return stocks, nil
}

func (r *StockRepository) DeleteByTicker(ticker string) error {
	return r.repo.DeleteByField("ticker", ticker, &StockEntity{})
}

func (r *StockRepository) FindBy(filters map[string]any) (*stock.Stock, error) {
	var entity StockEntity
	err := r.repo.FindOneBy(filters, &entity)
	if err != nil {
		return nil, err
	}
	return ToDomain(&entity), nil
}

func (r *StockRepository) FindPaginated(page int, limit int) ([]stock.Stock, int64, error) {
	var entities []StockEntity
	total, err := r.repo.FindPaginated(&entities, page, limit)
	if err != nil {
		return nil, 0, err
	}

	stocks := make([]stock.Stock, len(entities))
	for i, entity := range entities {
		domainStock := ToDomain(&entity)
		stocks[i] = *domainStock
	}
	return stocks, total, nil
}
