package api

import "time"

type PriceRecord struct {
	Source    string    `json:"source"`
	Currency  string    `json:"currency"`
	TimeStamp time.Time `json:"timeStamp"`
	Price     float64   `json:"price"`
}
