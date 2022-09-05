package prices

import (
	"context"
	"net/http"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/integrations/getgeoapi"
)

func GeoCurrencyProviderFunc(cfg *config.GeoCurrencyProvider) ProviderFunc {
	httpClient := http.DefaultClient
	httpClient.Timeout = time.Second * 5
	toCurrencies := make([]string, 0, len(cfg.Currencies))
	for _, val := range cfg.Currencies {
		toCurrencies = append(toCurrencies, val)
	}
	return func(ctx context.Context) ([]price.Value, error) {
		var values []price.Value
		for _, symbol := range cfg.Symbols {
			rates, err := getgeoapi.GetConversionRates(ctx, httpClient, cfg.ApiKey, symbol, 1, toCurrencies...)
			if err != nil {
				return nil, err
			}
			for key, rate := range rates.Rates {
				currency, found := cfg.Currencies[key]
				if !found {
					continue
				}
				values = append(values, price.Value{
					TimeStamp: rates.UpdateDate,
					Source:    "getgeoapi",
					Symbol:    symbol,
					Currency:  currency,
					Price:     rate.Rate,
				})
			}
		}
		return values, nil
	}
}
