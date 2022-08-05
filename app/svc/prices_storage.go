package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/price"

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
		if ps.mapping != nil {
			cur, found := ps.mapping[currency]
			if !found {
				continue
			}
			actualCurrency = cur
		}
		prices, err := ps.backend.LoadPrices(ctx, actualCurrency)
		if err != nil {
			continue
		}
		pricesGetter := price.NewAssetGetter(prices)
		for _, symbol := range symbols {
			value, err := pricesGetter.GetPrice(ctx, symbol, currency)
			if err != nil {
				continue
			}
			symbolPrices, found := result[symbol]
			if !found {
				symbolPrices = make(map[string]SymbolPrice)
			}
			symbolPrice := SymbolPrice{
				Source:    value.Source,
				Value:     value.Price,
				TimeStamp: value.TimeStamp,
			}
			if withHistory {
				history, err := pricesGetter.GetHistory(ctx, symbol, currency)
				if err == nil {
					historyList := make([]SymbolPricesHistory, 0, len(history))
					for _, historyItem := range history {
						historyList = append(historyList, SymbolPricesHistory{
							TimeStamp: historyItem.TimeStamp,
							Value:     historyItem.Price,
						})
					}
					symbolPrice.History = historyList
				}
			}
			symbolPrices[value.Currency] = symbolPrice
			result[symbol] = symbolPrices
		}
	}
	return result, nil
}
