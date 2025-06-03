package stock

import (
	"backend/domain/stock"
	"strconv"
	"time"
)

type StockEntity struct {
	ID     uint   `gorm:"column:id;primaryKey" json:"id"`
	Ticker string `gorm:"column:ticker;type:varchar(10)" json:"ticker"`
	// TargetFrom float64   `gorm:"column:target_from;type:decimal(10,2)" json:"target_from"`
	// TargetTo   float64   `gorm:"column:target_to;type:decimal(10,2)" json:"target_to"`
	Company    string    `gorm:"column:company;type:varchar(100)" json:"company"`
	Action     string    `gorm:"column:action;type:varchar(50)" json:"action"`
	Brokerage  string    `gorm:"column:brokerage;type:varchar(100)" json:"brokerage"`
	RatingFrom string    `gorm:"column:rating_from;type:varchar(20)" json:"rating_from"`
	RatingTo   string    `gorm:"column:rating_to;type:varchar(20)" json:"rating_to"`
	Time       time.Time `gorm:"column:time" json:"time"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	Category   string    `gorm:"column:category;type:varchar(50)" json:"category"`
}

func (StockEntity) TableName() string {
	return "stocks_ratings"
}

func ToDomain(e *StockEntity) *stock.Stock {
	if e == nil {
		return nil
	}
	return &stock.Stock{
		ID:     strconv.FormatUint(uint64(e.ID), 10),
		Ticker: e.Ticker,
		// TargetFrom: e.TargetFrom,
		// TargetTo:   e.TargetTo,
		Company:    e.Company,
		Action:     e.Action,
		Brokerage:  e.Brokerage,
		RatingFrom: e.RatingFrom,
		RatingTo:   e.RatingTo,
		Time:       e.Time,
		CreatedAt:  e.CreatedAt,
		Category:   e.Category,
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
		ID:     uint(id),
		Ticker: s.Ticker,
		// TargetFrom: s.TargetFrom,
		// TargetTo:   s.TargetTo,
		Company:    s.Company,
		Action:     s.Action,
		Brokerage:  s.Brokerage,
		RatingFrom: s.RatingFrom,
		RatingTo:   s.RatingTo,
		Time:       s.Time,
		CreatedAt:  s.CreatedAt,
		Category:   s.Category,
	}
}
