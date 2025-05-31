package stock_snapshot

import (
	"backend/domain/stock_snapshot"
)

type StockSnapshotService struct {
	repo stock_snapshot.Repository
}

func NewStockSnapshotService(repo stock_snapshot.Repository) *StockSnapshotService {
	return &StockSnapshotService{repo: repo}
}

func (s *StockSnapshotService) FindPaginated(page, limit int) ([]stock_snapshot.StockSnapshot, int64, error) {
	return s.repo.FindPaginated(page, limit)
}

func (s *StockSnapshotService) FindByID(id string) (*stock_snapshot.StockSnapshot, error) {
	filters := map[string]any{"id": id}
	return s.repo.FindBy(filters)
}
