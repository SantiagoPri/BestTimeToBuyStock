package stock_snapshot

import (
	"backend/domain/stock_snapshot"
	stockRepo "backend/infrastructure/repositories/stock"
	"strconv"
	"time"
)

type StockSnapshotEntity struct {
	ID          uint                  `gorm:"column:id;primaryKey" json:"id"`
	StockID     uint                  `gorm:"column:stock_id;not null;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"stock_id"`
	Week        uint                  `gorm:"column:week;not null" json:"week"`
	RatingFrom  string                `gorm:"column:rating_from;type:varchar(20)" json:"rating_from"`
	RatingTo    string                `gorm:"column:rating_to;type:varchar(20)" json:"rating_to"`
	TargetFrom  float64               `gorm:"column:target_from;type:decimal(10,2)" json:"target_from"`
	TargetTo    float64               `gorm:"column:target_to;type:decimal(10,2)" json:"target_to"`
	Price       float64               `gorm:"column:price;type:decimal(10,2)" json:"price"`
	Action      string                `gorm:"column:action;type:varchar(50)" json:"action"`
	NewsTitle   string                `gorm:"column:news_title;type:varchar(200)" json:"news_title"`
	NewsSummary string                `gorm:"column:news_summary;type:text" json:"news_summary"`
	CreatedAt   time.Time             `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Stock       stockRepo.StockEntity `gorm:"foreignKey:StockID;references:ID" json:"stock"`
}

func (StockSnapshotEntity) TableName() string {
	return "stock_snapshots"
}

func ToDomain(e *StockSnapshotEntity) *stock_snapshot.StockSnapshot {
	if e == nil {
		return nil
	}
	return &stock_snapshot.StockSnapshot{
		ID:          strconv.FormatUint(uint64(e.ID), 10),
		Week:        e.Week,
		RatingFrom:  e.RatingFrom,
		RatingTo:    e.RatingTo,
		TargetFrom:  e.TargetFrom,
		TargetTo:    e.TargetTo,
		Price:       e.Price,
		Action:      e.Action,
		NewsTitle:   e.NewsTitle,
		NewsSummary: e.NewsSummary,
		CreatedAt:   e.CreatedAt,
		Stock:       *stockRepo.ToDomain(&e.Stock),
	}
}

func FromDomain(s *stock_snapshot.StockSnapshot) *StockSnapshotEntity {
	if s == nil {
		return nil
	}
	id, err := strconv.ParseUint(s.ID, 10, 64)
	if err != nil {
		id = 0 // For create operations
	}
	stockEntity := stockRepo.FromDomain(&s.Stock)
	return &StockSnapshotEntity{
		ID:          uint(id),
		StockID:     stockEntity.ID,
		Week:        s.Week,
		RatingFrom:  s.RatingFrom,
		RatingTo:    s.RatingTo,
		TargetFrom:  s.TargetFrom,
		TargetTo:    s.TargetTo,
		Price:       s.Price,
		Action:      s.Action,
		NewsTitle:   s.NewsTitle,
		NewsSummary: s.NewsSummary,
		CreatedAt:   s.CreatedAt,
		Stock:       *stockEntity,
	}
}
