package stock_snapshot

import (
	"backend/domain/stock_snapshot"
	"backend/infrastructure/repositories"

	"gorm.io/gorm"
)

// Compile-time assertion to ensure StockSnapshotRepository implements stock_snapshot.Repository
var _ stock_snapshot.Repository = (*StockSnapshotRepository)(nil)

type StockSnapshotRepository struct {
	db   *gorm.DB
	repo repositories.Repository
}

func NewStockSnapshotRepository(db *gorm.DB) *StockSnapshotRepository {
	return &StockSnapshotRepository{
		db:   db,
		repo: repositories.NewBaseRepository(db, &StockSnapshotEntity{}),
	}
}

func (r *StockSnapshotRepository) Save(s *stock_snapshot.StockSnapshot) error {
	entity := FromDomain(s)
	return r.repo.Save(entity)
}

func (r *StockSnapshotRepository) FindAll() ([]stock_snapshot.StockSnapshot, error) {
	var entities []StockSnapshotEntity
	err := r.db.Preload("Stock").Find(&entities).Error
	if err != nil {
		return nil, err
	}

	snapshots := make([]stock_snapshot.StockSnapshot, len(entities))
	for i, entity := range entities {
		domainSnapshot := ToDomain(&entity)
		snapshots[i] = *domainSnapshot
	}
	return snapshots, nil
}

func (r *StockSnapshotRepository) FindBy(filters map[string]any) (*stock_snapshot.StockSnapshot, error) {
	var entity StockSnapshotEntity
	convertedFilters := repositories.ParseIDFilter(filters, "id", "stock_id")
	err := r.db.Preload("Stock").Where(convertedFilters).First(&entity).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, stock_snapshot.ErrNotFound
		}
		return nil, err
	}
	return ToDomain(&entity), nil
}

func (r *StockSnapshotRepository) FindByCategory(category string) ([]stock_snapshot.StockSnapshot, error) {
	var entities []StockSnapshotEntity
	err := r.db.Preload("Stock").
		Joins("JOIN stocks ON stock_snapshots.stock_id = stocks.id").
		Where("stocks.category = ?", category).
		Find(&entities).Error
	if err != nil {
		return nil, err
	}

	snapshots := make([]stock_snapshot.StockSnapshot, len(entities))
	for i, entity := range entities {
		domainSnapshot := ToDomain(&entity)
		snapshots[i] = *domainSnapshot
	}
	return snapshots, nil
}

func (r *StockSnapshotRepository) FindPaginated(page int, limit int) ([]stock_snapshot.StockSnapshot, int64, error) {
	var entities []StockSnapshotEntity
	var total int64

	err := r.db.Model(&StockSnapshotEntity{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err = r.db.Preload("Stock").
		Offset(offset).
		Limit(limit).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	snapshots := make([]stock_snapshot.StockSnapshot, len(entities))
	for i, entity := range entities {
		domainSnapshot := ToDomain(&entity)
		snapshots[i] = *domainSnapshot
	}
	return snapshots, total, nil
}
