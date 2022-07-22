package api

import (
	"strconv"
	"strings"
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
	strVal := strings.TrimRight(strconv.FormatFloat(float64(d), 'f', 10, 64), "0")
	return []byte(strings.TrimRight(strVal, ".")), nil
}
