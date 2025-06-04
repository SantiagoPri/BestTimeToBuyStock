package stock

import (
	"context"
	"fmt"
	"strings"

	"backend/domain/stock"
	"backend/infrastructure/repositories"
	"backend/pkg/errors"

	"gorm.io/gorm"
)

var _ stock.Repository = (*StockRepository)(nil)

type StockRepository struct {
	repo repositories.Repository
	db   *gorm.DB
}

func NewStockRepository(db *gorm.DB) *StockRepository {
	return &StockRepository{
		repo: repositories.NewBaseRepository(db, &StockEntity{}),
		db:   db,
	}
}

func (r *StockRepository) Save(s *stock.Stock) error {
	entity := FromDomain(s)
	if err := r.repo.Save(entity); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to save stock", err)
	}
	return nil
}

func (r *StockRepository) FindAll() ([]stock.Stock, error) {
	var entities []StockEntity
	err := r.repo.FindAll(&entities)
	if err != nil {
		return nil, errors.Wrap(errors.ErrInternal, "failed to find stocks", err)
	}

	stocks := make([]stock.Stock, len(entities))
	for i, entity := range entities {
		domainStock := ToDomain(&entity)
		stocks[i] = *domainStock
	}
	return stocks, nil
}

func (r *StockRepository) DeleteByTicker(ticker string) error {
	if err := r.repo.DeleteByField("ticker", ticker, &StockEntity{}); err != nil {
		return errors.Wrap(errors.ErrInternal, "failed to delete stock", err)
	}
	return nil
}

func (r *StockRepository) FindBy(filters map[string]any) (*stock.Stock, error) {
	var entity StockEntity
	err := r.repo.FindOneBy(filters, &entity)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, stock.ErrNotFound
		}
		return nil, errors.Wrap(errors.ErrInternal, "failed to find stock", err)
	}
	return ToDomain(&entity), nil
}

func (r *StockRepository) FindPaginated(page int, limit int) ([]stock.Stock, int64, error) {
	var entities []StockEntity
	total, err := r.repo.FindPaginated(&entities, page, limit)
	if err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to find paginated stocks", err)
	}

	stocks := make([]stock.Stock, len(entities))
	for i, entity := range entities {
		domainStock := ToDomain(&entity)
		stocks[i] = *domainStock
	}
	return stocks, total, nil
}

func (r *StockRepository) PickStocksForSession(categories []string) ([]stock.Stock, error) {
	if len(categories) != 3 {
		return nil, errors.New(errors.ErrInvalidInput, "exactly 3 categories required")
	}

	result := make([]stock.Stock, 0, 12)

	for _, category := range categories {
		var entities []StockEntity
		err := r.repo.FindRandomByField("category", category, 4, &entities)
		if err != nil {
			return nil, errors.Wrap(errors.ErrInternal, "failed to fetch stocks for category", err)
		}

		for _, entity := range entities {
			domainStock := ToDomain(&entity)
			result = append(result, *domainStock)
		}
	}

	return result, nil
}

func (r *StockRepository) applyFilters(query *gorm.DB, filters map[string]string) (*gorm.DB, error) {
	for field, value := range filters {
		filterType, exists := stock.ValidFilters[field]
		if !exists {
			return nil, errors.New(errors.ErrInvalidInput, fmt.Sprintf("invalid filter field: %s", field))
		}

		switch filterType {
		case stock.FilterExact:
			query = query.Where(fmt.Sprintf("%s = ?", field), value)
		case stock.FilterILike:
			query = query.Where(fmt.Sprintf("LOWER(%s) LIKE ?", field), "%"+strings.ToLower(value)+"%")
		}
	}
	return query, nil
}

func (r *StockRepository) applySorting(query *gorm.DB, sortBy, order string) *gorm.DB {
	if sortBy != "" {
		direction := "ASC"
		if strings.ToLower(order) == "desc" {
			direction = "DESC"
		}
		query = query.Order(fmt.Sprintf("%s %s", sortBy, direction))
	}
	return query
}

func (r *StockRepository) FindAllStocks(ctx context.Context, params stock.QueryParams) ([]stock.Stock, int64, error) {
	var entities []StockEntity
	var total int64

	query := r.db.Model(&StockEntity{}).WithContext(ctx)

	var err error
	query, err = r.applyFilters(query, params.Filters)
	if err != nil {
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to count stocks", err)
	}

	query = r.applySorting(query, params.SortBy, params.SortOrder)

	offset := (params.Page - 1) * params.PageSize
	query = query.Offset(offset).Limit(params.PageSize)

	if err := query.Find(&entities).Error; err != nil {
		return nil, 0, errors.Wrap(errors.ErrInternal, "failed to fetch stocks", err)
	}

	stocks := make([]stock.Stock, len(entities))
	for i, entity := range entities {
		domainStock := ToDomain(&entity)
		stocks[i] = *domainStock
	}

	return stocks, total, nil
}
