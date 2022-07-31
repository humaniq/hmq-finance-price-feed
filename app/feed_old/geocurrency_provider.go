package feed_old

import (
	"context"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/prices_old"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type GeoCurrencyPriceProvider struct {
	client     *prices_old.IPCurrencyAPI
	ticker     *time.Ticker
	symbols    map[string]string
	currencies map[string]string
}

func NewGeoCurrencyPriceProvider(tick time.Duration, client *prices_old.IPCurrencyAPI, symbols map[string]string, currencies map[string]string) *GeoCurrencyPriceProvider {
	return &GeoCurrencyPriceProvider{
		client:     client,
		ticker:     time.NewTicker(tick),
		symbols:    symbols,
		currencies: currencies,
	}
}

func (gcpp *GeoCurrencyPriceProvider) Provide(ctx context.Context, feed chan<- []price.Value) error {
	symbolsList := make([]string, 0, len(gcpp.symbols))
	for symbol, _ := range gcpp.symbols {
		symbolsList = append(symbolsList, symbol)
	}
	logger.Info(ctx, "providing geocurrency: %+v for %+v", symbolsList, gcpp.currencies)
	for range gcpp.ticker.C {
		now := time.Now()
		logger.Info(ctx, "geoCurrency tick: %s", now)
		var items []price.Value
		for currency, currencyKey := range gcpp.currencies {
			response, err := gcpp.client.GetConversionRates(ctx, currency, 1, symbolsList...)
			if err != nil {
				logger.Error(ctx, "error getting from geoCurrency: %s", err.Error())
				continue
			}
			for key, value := range response {
				items = append(items, price.Value{
					TimeStamp: now,
					Source:    "geo",
					Symbol:    gcpp.symbols[key],
					Currency:  currencyKey,
					Price:     value,
				})
			}
		}
		feed <- items
	}
	return nil
}
