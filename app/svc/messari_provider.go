package svc

import (
	"context"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/prices_old"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type MessariPriceProvider struct {
	client  *prices_old.Messari
	ticker  *time.Ticker
	symbols []string
}

func NewMessariPriceProvider(tick time.Duration, symbols []string) *MessariPriceProvider {
	return &MessariPriceProvider{
		client:  prices_old.NewMessari(),
		ticker:  time.NewTicker(tick),
		symbols: symbols,
	}
}
func (mpp *MessariPriceProvider) Provide(ctx context.Context, feed chan<- *FeedItem) error {
	for range mpp.ticker.C {
		feedItem := &FeedItem{records: make([]*PriceRecord, 0, len(mpp.symbols)*3)}
		for _, symbol := range mpp.symbols {
			price, err := mpp.client.GetPricesForSymbol(ctx, symbol)
			if err != nil {
				logger.Error(ctx, "error getting prices_old: %s", err.Error())
				continue
			}
			now := time.Now()
			feedItem.records = append(feedItem.records, &PriceRecord{
				Source:    "messari",
				Symbol:    symbol,
				Currency:  "ETH",
				Price:     price.Eth,
				TimeStamp: now,
			})
			feedItem.records = append(feedItem.records, &PriceRecord{
				Source:    "messari",
				Symbol:    symbol,
				Currency:  "BTC",
				Price:     price.Btc,
				TimeStamp: now,
			})
			feedItem.records = append(feedItem.records, &PriceRecord{
				Source:    "messari",
				Symbol:    symbol,
				Currency:  "USD",
				Price:     price.Usd,
				TimeStamp: now,
			})
		}
		logger.Info(ctx, "MESSARI: %+v", feedItem)
		feed <- feedItem
	}
	return nil
}
func (mpp *MessariPriceProvider) Stop() {
	mpp.ticker.Stop()
}
