package price

import (
	"context"
	"sort"

	"github.com/humaniq/hmq-finance-price-feed/app"
)

type AssetGetter struct {
	asset *Asset
}

func NewAssetGetter(asset *Asset) *AssetGetter {
	return &AssetGetter{asset: asset}
}
func (ag *AssetGetter) GetPrice(ctx context.Context, symbol string, currency string) (*Value, error) {
	symbolPriceValue, found := ag.asset.Prices[symbol]
	if !found {
		return nil, app.ErrNotFound
	}
	if symbolPriceValue.Price == 0 {
		return nil, app.ErrValueInvalid
	}
	if currency != ag.asset.Name {
		currencyPriceValue, found := ag.asset.Prices[currency]
		if !found {
			return nil, app.ErrNotFound
		}
		if currencyPriceValue.Price == 0 {
			return nil, app.ErrValueInvalid
		}
		return &Value{
			TimeStamp: symbolPriceValue.TimeStamp,
			Source:    "estimation",
			Symbol:    symbol,
			Currency:  currency,
			Price:     symbolPriceValue.Price / currencyPriceValue.Price,
		}, nil
	}
	return &symbolPriceValue, nil
}
func (ag *AssetGetter) GetHistory(ctx context.Context, symbol string, currency string) (History, error) {
	symbolHistory, found := ag.asset.History[symbol]
	if !found || len(symbolHistory) == 0 {
		return nil, app.ErrNotFound
	}
	if currency != ag.asset.Name {
		currencyHistory, found := ag.asset.History[currency]
		if !found || len(currencyHistory) == 0 {
			return nil, app.ErrNotFound
		}
		history := make([]struct {
			isSymbol bool
			record   HistoryRecord
		}, 0, len(symbolHistory)+len(currencyHistory))
		for _, symbolHistoryRecord := range symbolHistory {
			history = append(history, struct {
				isSymbol bool
				record   HistoryRecord
			}{isSymbol: true, record: symbolHistoryRecord})
		}
		for _, currencyHistoryRecord := range currencyHistory {
			history = append(history, struct {
				isSymbol bool
				record   HistoryRecord
			}{isSymbol: false, record: currencyHistoryRecord})
		}
		sort.Slice(history, func(i, j int) bool {
			return history[i].record.TimeStamp.Before(history[j].record.TimeStamp)
		})
		priceHistory := make(History, 0, len(history))
		iterationSymbolPrice := float64(0)
		iterationCurrencyPrice := float64(0)
		for _, record := range history {
			if record.isSymbol {
				iterationSymbolPrice = record.record.Price
			} else {
				iterationCurrencyPrice = record.record.Price
			}
			if iterationSymbolPrice != 0 && iterationCurrencyPrice != 0 {
				priceHistory = append(priceHistory, HistoryRecord{
					TimeStamp: record.record.TimeStamp,
					Price:     iterationSymbolPrice / iterationCurrencyPrice,
				})
			}
		}
		if len(priceHistory) == 0 {
			return nil, app.ErrNotFound
		}
		return priceHistory, nil
	}
	return symbolHistory, nil
}
