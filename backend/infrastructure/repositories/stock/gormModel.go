package stock

import (
	"backend/domain/stock"
	"time"
)

type StockEntity struct {
	Ticker     string    `gorm:"primaryKey;type:varchar(10)"`
	TargetFrom string    `gorm:"type:varchar(10)"`
	TargetTo   string    `gorm:"type:varchar(10)"`
	Company    string    `gorm:"type:varchar(100)"`
	Action     string    `gorm:"type:varchar(20)"`
	Brokerage  string    `gorm:"type:varchar(50)"`
	RatingFrom string    `gorm:"type:varchar(20)"`
	RatingTo   string    `gorm:"type:varchar(20)"`
	Time       time.Time `gorm:"type:timestamp"`
}

func (StockEntity) TableName() string {
	return "stock_ratings"
}

func ToDomain(e *StockEntity) *stock.Stock {
	if e == nil {
		return nil
	}
	return &stock.Stock{
		Ticker:     e.Ticker,
		TargetFrom: e.TargetFrom,
		TargetTo:   e.TargetTo,
		Company:    e.Company,
		Action:     e.Action,
		Brokerage:  e.Brokerage,
		RatingFrom: e.RatingFrom,
		RatingTo:   e.RatingTo,
		Time:       e.Time,
	}
}

func FromDomain(s *stock.Stock) *StockEntity {
	if s == nil {
		return nil
	}
	return &StockEntity{
		Ticker:     s.Ticker,
		TargetFrom: s.TargetFrom,
		TargetTo:   s.TargetTo,
		Company:    s.Company,
		Action:     s.Action,
		Brokerage:  s.Brokerage,
		RatingFrom: s.RatingFrom,
		RatingTo:   s.RatingTo,
		Time:       s.Time,
	}
}
