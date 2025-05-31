package stock

import "time"

type Stock struct {
	ID        string    `json:"id"`
	Ticker    string    `json:"ticker"`
	Company   string    `json:"company"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
}
