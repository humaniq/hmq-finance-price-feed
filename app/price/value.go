package price

import "time"

type Value struct {
	TimeStamp time.Time `json:"time"`
	Source    string    `json:"source"`
	Symbol    string    `json:"symbol"`
	Currency  string    `json:"currency"`
	Price     float64   `json:"price"`
}
