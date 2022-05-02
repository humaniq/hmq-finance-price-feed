package service

import (
	"context"
	"time"
)

type StorageRecord struct {
	Source    string             `datastore:"source"`
	TimeStamp time.Time          `datastore:"timeStamp"`
	Prices    map[string]float64 `datastore:"prices"`
}

type PriceRecord struct {
	Source        string    `json:"source"`
	Symbol        string    `json:"symbol"`
	Currency      string    `json:"currency"`
	Price         float64   `json:"price"`
	PreviousPrice float64   `json:"previousPrice,omitempty"`
	TimeStamp     time.Time `json:"timeStamp"`
}

func NewPriceRecord(symbol string, currency string, price float64, source string) *PriceRecord {
	return &PriceRecord{
		Source:    source,
		Symbol:    symbol,
		Currency:  currency,
		Price:     price,
		TimeStamp: time.Now(),
	}
}
func (pr *PriceRecord) WithPreviousPrice(price float64) *PriceRecord {
	pr.PreviousPrice = price
	return pr
}

type Prices map[string]float64

type PricesSetter interface {
	SetPrices(ctx context.Context, source string, timeStamp time.Time, prices Prices) error
}
type PricesGetter interface {
	GetLatestPrices(ctx context.Context, symbols []string, currencies []string) (map[string]Prices, error)
}
type Pricer interface {
	PricesGetter
	PricesSetter
}
