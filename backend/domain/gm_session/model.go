package gm_session

import (
	"backend/domain/stock"
	"context"
)

type AI interface {
	GetGMResponse(ctx context.Context, categories []string, stocks []stock.Stock) (map[string]*GMWeekData, error)
}

type GMWeekData struct {
	Headlines []string           `json:"headlines"`
	Stocks    []StockWeekInsight `json:"stocks"`
}

type StockWeekInsight struct {
	Ticker      string  `json:"ticker"`
	CompanyName string  `json:"company_name"`
	RatingFrom  string  `json:"rating_from"`
	RatingTo    string  `json:"rating_to"`
	Action      string  `json:"action"`
	Price       float64 `json:"price"`
	PriceChange float64 `json:"price_change"`
}
