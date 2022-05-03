package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/prices"
)

type CoinGeckoPriceProvider struct {
	get    func(ctx context.Context) (map[string]map[string]float64, error)
	ticker *time.Ticker
}

func NewCoinGeckoProvider(tick time.Duration, client *prices.CoinGecko, symbols []string, currencies []string) *CoinGeckoPriceProvider {
	return &CoinGeckoPriceProvider{get: client.GetterFunc(symbols, currencies), ticker: time.NewTicker(tick)}
}

func (cgpp *CoinGeckoPriceProvider) Provide(ctx context.Context, feed chan<- *FeedItem) error {
	for range cgpp.ticker.C {
		result, err := cgpp.get(ctx)
		if err != nil {
			logger.Error(ctx, "error getting from coingecko: %s", err.Error())
			continue
		}
		now := time.Now()
		var item FeedItem
		for key, value := range result {
			for k, v := range value {
				item.records = append(item.records, &PriceRecord{
					Source:    "coingecko",
					Symbol:    key,
					Currency:  k,
					Price:     v,
					TimeStamp: now,
				})
			}
		}
		feed <- &item
	}
	return nil
}