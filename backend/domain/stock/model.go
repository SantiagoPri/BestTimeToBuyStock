package stock

import "time"

type Stock struct {
	ID     string `json:"id"`
	Ticker string `json:"ticker"`
	// TargetFrom float64   `json:"target_from"`
	// TargetTo   float64   `json:"target_to"`
	Company    string    `json:"company"`
	Action     string    `json:"action"`
	Brokerage  string    `json:"brokerage"`
	RatingFrom string    `json:"rating_from"`
	RatingTo   string    `json:"rating_to"`
	Time       time.Time `json:"time"`
	CreatedAt  time.Time `json:"created_at"`
	Category   string    `json:"category"`
}
