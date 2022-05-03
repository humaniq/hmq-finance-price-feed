package svc

import (
	"context"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
)

type PricesStorage struct {
	backend storage.SymbolPricesGetter
}

func NewPricesStorage(backend storage.SymbolPricesGetter) *PricesStorage {
	return &PricesStorage{backend: backend}
}

func (ps *PricesStorage) GetPrices(ctx context.Context, symbols []string, currencies []string) (map[string]SymbolPrices, error) {
	result := make(map[string]SymbolPrices)
	for _, symbol := range symbols {
		symbolPrices, err := ps.backend.GetSymbolPrices(ctx, symbol)
		if err != nil {
			if !errors.Is(err, storage.ErrNotFound) {
				return nil, err
			}
			continue
		}
		prices := make(map[string]SymbolPrice)
		for currency, price := range symbolPrices.Prices {
			prices[currency] = SymbolPrice{
				Source:    symbolPrices.Source,
				Value:     price,
				TimeStamp: symbolPrices.TimeStamp,
			}
		}
		result[symbol] = prices
	}
	return result, nil
}
