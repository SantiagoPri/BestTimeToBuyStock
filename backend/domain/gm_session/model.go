package gm_session

type GMWeekData struct {
	Headlines []string           `json:"headlines"`
	Stocks    []StockWeekInsight `json:"stocks"`
}

type StockWeekInsight struct {
	Ticker     string  `json:"ticker"`
	RatingFrom string  `json:"rating_from"`
	RatingTo   string  `json:"rating_to"`
	Action     string  `json:"action"`
	Price      float64 `json:"price"`
}
