package stock_snapshot

import (
	"backend/domain/stock"
	"time"
)

type StockSnapshot struct {
	ID          uint
	Week        uint
	RatingFrom  string
	RatingTo    string
	TargetFrom  float64
	TargetTo    float64
	Price       float64
	Action      string
	NewsTitle   string
	NewsSummary string
	CreatedAt   time.Time
	Stock       stock.Stock
}
