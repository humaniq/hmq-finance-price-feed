package feed

import (
	"context"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type GeoCurrencyPriceProvider struct {
	client     *prices.IPCurrencyAPI
	ticker     *time.Ticker
	symbols    map[string]string
	currencies map[string]string
}

func NewGeoCurrencyPriceProvider(tick time.Duration, client *prices.IPCurrencyAPI, symbols map[string]string, currencies map[string]string) *GeoCurrencyPriceProvider {
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
		var items []price.Value
		for currency, currencyKey := range gcpp.currencies {
			response, err := gcpp.client.GetConversionRates(ctx, currency, 1, symbolsList...)
			if err != nil {
				logger.Error(ctx, "error getting from geoCurrency: %s", err.Error())
				continue
			}
			now := time.Now()
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
