package api

import "time"

type PriceRecord struct {
	Source    string                `json:"source"`
	Currency  string                `json:"currency"`
	TimeStamp time.Time             `json:"time"`
	Price     float64               `json:"price"`
	History   []*PriceHistoryRecord `json:"history,omitempty"`
}
type PriceHistoryRecord struct {
	TimeStamp time.Time `json:"time"`
	Price     float64   `json:"price"`
}
