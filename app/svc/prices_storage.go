package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
)

type Prices struct {
	backend  storage.PricesLoader
	required []string
	std      string
	mapping  map[string]string
}

func NewPrices(backend storage.PricesLoader) *Prices {
	return &Prices{backend: backend}
}
func (ps *Prices) WithMapping(mapping map[string]string) *Prices {
	ps.mapping = mapping
	return ps
}

func (ps *Prices) GetPrices(ctx context.Context, symbols []string, currencies []string, withHistory bool) (map[string]SymbolPrices, error) {
	result := make(map[string]SymbolPrices)
	for _, currency := range currencies {
		actualCurrency := currency
		estimate := false
		if ps.mapping != nil {
			cur, found := ps.mapping[currency]
			if !found {
				continue
			}
			if cur != currency {
				estimate = true
			}
			actualCurrency = cur
		}
		prices, err := ps.backend.LoadPrices(ctx, actualCurrency)
		if err != nil {
			return nil, err
		}
		for _, symbol := range symbols {
			var value *state.Price
			if estimate {
				value = prices.Estimate(symbol, currency, withHistory)
			} else {
				value = prices.Get(symbol, withHistory)
			}
			if value == nil {
				continue
			}
			record, found := result[symbol]
			if !found {
				record = make(map[string]SymbolPrice)
			}
			val := SymbolPrice{
				Source:    value.Current.Source,
				Value:     value.Current.Price,
				TimeStamp: value.Current.TimeStamp,
			}
			if withHistory {
				for _, rec := range value.History {
					val.History = append(val.History, SymbolPricesHistory{TimeStamp: rec.TimeStamp, Value: rec.Value})
				}
			}
			record[currency] = val
			result[symbol] = record
		}
	}
	return result, nil
}
