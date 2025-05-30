package stock

import (
	"backend/domain/stock"
	"time"
)

type StockEntity struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Ticker    string    `gorm:"type:varchar(10)" json:"ticker"`
	Company   string    `gorm:"type:varchar(100)" json:"company"`
	Category  string    `gorm:"type:varchar(50)" json:"category"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (StockEntity) TableName() string {
	return "stocks"
}

func ToDomain(e *StockEntity) *stock.Stock {
	if e == nil {
		return nil
	}
	return &stock.Stock{
		ID:        e.ID,
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
	return &StockEntity{
		ID:        s.ID,
		Ticker:    s.Ticker,
		Company:   s.Company,
		Category:  s.Category,
		CreatedAt: s.CreatedAt,
	}
}
