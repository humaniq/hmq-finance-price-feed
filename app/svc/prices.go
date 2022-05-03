package svc

import (
	"context"
	"time"
)

type SymbolPrices map[string]SymbolPrice
type SymbolPrice struct {
	Source    string    `json:"source"`
	Value     float64   `json:"value"`
	TimeStamp time.Time `json:"timeStamp"`
}

type PricesGetter interface {
	GetPrices(ctx context.Context, symbols []string, currencies []string) (map[string]SymbolPrices, error)
}
