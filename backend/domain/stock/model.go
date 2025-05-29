package stock

import "time"

type Stock struct {
	Ticker     string
	TargetFrom string
	TargetTo   string
	Company    string
	Action     string
	Brokerage  string
	RatingFrom string
	RatingTo   string
	Time       time.Time
}
