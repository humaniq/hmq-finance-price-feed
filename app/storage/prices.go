package storage

import (
	"context"
	"time"
)

type PricesRecord struct {
	Symbol    string
	Source    string
	TimeStamp time.Time
	Prices    map[string]float64
}

func NewPricesRecord(symbol string, source string, timeStamp time.Time) *PricesRecord {
	return &PricesRecord{
		Symbol:    symbol,
		Source:    source,
		TimeStamp: timeStamp,
		Prices:    make(map[string]float64),
	}
}

type SymbolPricesSetter interface {
	SetSymbolPrices(ctx context.Context, symbol string, source string, timeStamp time.Time, prices map[string]float64) error
}
type SymbolPricesGetter interface {
	GetSymbolPrices(ctx context.Context, symbol string) (*PricesRecord, error)
}
type SymbolPrices interface {
	SymbolPricesSetter
	SymbolPricesGetter
}
