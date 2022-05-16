package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
)

type Prices struct {
	backend  storage.PricesLoader
	required []string
	std      string
}

func NewPrices(backend storage.PricesLoader) *Prices {
	return &Prices{backend: backend}
}
func (ps *Prices) WithStdCurrency(std string) *Prices {
	ps.std = std
	return ps
}

func (ps *Prices) GetPrices(ctx context.Context, symbols []string, currencies []string) (map[string]SymbolPrices, error) {
	result := make(map[string]SymbolPrices)
	for _, currency := range currencies {
		prices, err := ps.backend.LoadPrices(ctx, currency)
		if err != nil {
			return nil, err
		}
		for _, symbol := range symbols {
			value, found := prices.Values()[symbol]
			if !found {
				continue
			}
			record, found := result[symbol]
			if !found {
				record = make(map[string]SymbolPrice)
			}
			record[currency] = SymbolPrice{
				Source:    value.Source,
				Value:     value.Price,
				TimeStamp: value.TimeStamp,
			}
		}
	}

	return result, nil
}
