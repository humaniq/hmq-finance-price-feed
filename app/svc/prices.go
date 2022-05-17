package svc

import (
	"context"
	"time"
)

type SymbolPrices map[string]SymbolPrice
type SymbolPrice struct {
	Source    string                `json:"source"`
	Value     float64               `json:"value"`
	TimeStamp time.Time             `json:"time"`
	History   []SymbolPricesHistory `json:"history,omitempty"`
}
type SymbolPricesHistory struct {
	TimeStamp time.Time `json:"time"`
	Value     float64   `json:"value"`
}

type PricesGetter interface {
	GetPrices(ctx context.Context, symbols []string, currencies []string, withHistory bool) (map[string]SymbolPrices, error)
}
