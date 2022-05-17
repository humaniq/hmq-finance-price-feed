package feed

import (
	"context"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type CoinGeckoPriceProvider struct {
	get    func(ctx context.Context) (map[string]map[string]float64, error)
	ticker *time.Ticker
}

func NewCoinGeckoProvider(tick time.Duration, client *prices.CoinGecko, symbols map[string]string, currencies map[string]string) *CoinGeckoPriceProvider {
	return &CoinGeckoPriceProvider{get: client.GetterFunc(symbols, currencies), ticker: time.NewTicker(tick)}
}

func (cgpp *CoinGeckoPriceProvider) Provide(ctx context.Context, feed chan<- []price.Value) error {
	for range cgpp.ticker.C {
		result, err := cgpp.get(ctx)
		if err != nil {
			logger.Error(ctx, "error getting from coingecko: %s", err.Error())
			continue
		}
		now := time.Now()
		var items []price.Value
		for key, value := range result {
			for k, v := range value {
				items = append(items, price.Value{
					TimeStamp: now,
					Source:    "coingecko",
					Symbol:    key,
					Currency:  k,
					Price:     v,
				})
			}
		}
		feed <- items
	}
	return nil
}
