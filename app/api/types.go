package api

import (
	"strconv"
	"time"
)

type PriceRecord struct {
	Source    string                `json:"source"`
	Currency  string                `json:"currency"`
	TimeStamp time.Time             `json:"time"`
	Price     Decimal               `json:"price"`
	History   []*PriceHistoryRecord `json:"history,omitempty"`
}
type PriceHistoryRecord struct {
	TimeStamp time.Time `json:"time"`
	Price     Decimal   `json:"price"`
}

type Decimal float64

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatFloat(float64(d), 'f', -1, 64)), nil
}
