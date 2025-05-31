package stock

import (
	"backend/domain/stock"
	"strconv"
	"time"
)

type StockEntity struct {
	ID        uint      `gorm:"column:id;primaryKey" json:"id"`
	Ticker    string    `gorm:"column:ticker;type:varchar(10)" json:"ticker"`
	Company   string    `gorm:"column:company;type:varchar(100)" json:"company"`
	Category  string    `gorm:"column:category;type:varchar(50)" json:"category"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (StockEntity) TableName() string {
	return "stocks"
}

func ToDomain(e *StockEntity) *stock.Stock {
	if e == nil {
		return nil
	}
	return &stock.Stock{
		ID:        strconv.FormatUint(uint64(e.ID), 10),
		Ticker:    e.Ticker,
		Company:   e.Company,
		Category:  e.Category,
		CreatedAt: e.CreatedAt,
	}
}

func FromDomain(s *stock.Stock) *StockEntity {
	if s == nil {
		return nil
	}
	id, err := strconv.ParseUint(s.ID, 10, 64)
	if err != nil {
		id = 0 // For create operations
	}
	return &StockEntity{
		ID:        uint(id),
		Ticker:    s.Ticker,
		Company:   s.Company,
		Category:  s.Category,
		CreatedAt: s.CreatedAt,
	}
}
