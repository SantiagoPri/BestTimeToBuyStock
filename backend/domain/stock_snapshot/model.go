package stock_snapshot

import (
	"backend/domain/stock"
	"time"
)

type StockSnapshot struct {
	ID          string      `json:"id"`
	Week        uint        `json:"week"`
	RatingFrom  string      `json:"ratingFrom"`
	RatingTo    string      `json:"ratingTo"`
	TargetFrom  float64     `json:"targetFrom"`
	TargetTo    float64     `json:"targetTo"`
	Price       float64     `json:"price"`
	Action      string      `json:"action"`
	NewsTitle   string      `json:"newsTitle"`
	NewsSummary string      `json:"newsSummary"`
	CreatedAt   time.Time   `json:"createdAt"`
	Stock       stock.Stock `json:"stock"`
}
