package svc

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type PriceFeedConsumerFunc func() chan PriceRecord

func ConsumerForDS(ps *PriceSvc) chan PriceRecord {
	ctx := context.Background()
	queue := make(chan PriceRecord)
	go func() {
		for record := range queue {
			if err := ps.SetSymbolPrice(ctx, record.Symbol, record.Currency, record.Price, record.Source); err != nil {
				logger.Error(ctx, "[DS] error setting symbol price: %s", err)
			}
		}
	}()
	return queue
}
func ConsumerForLog() chan PriceRecord {
	ctx := context.Background()
	queue := make(chan PriceRecord)
	go func() {
		for record := range queue {
			logger.Info(ctx, "got record: %+v", record)
		}
	}()
	return queue
}
