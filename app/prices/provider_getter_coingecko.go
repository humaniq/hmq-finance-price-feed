package prices

import (
	"context"
	"net/http"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/integrations/coingecko"
)

type CoinGecko struct {
	symbolMapper   map[string]string
	currencyMapper map[string]string
}

func NewCoingecko(assets *config.CoinGeckoAssets) *CoinGecko {
	return &CoinGecko{
		symbolMapper:   assets.Symbols,
		currencyMapper: assets.Currencies,
	}
}

func (cg *CoinGecko) GetterFunc(symbols []string, currencies []string) ProviderFunc {
	smbTokens := make([]string, 0, len(symbols))
	smbMapper := make(map[string]string)
	for _, symbol := range symbols {
		token, found := cg.symbolMapper[symbol]
		if found {
			smbTokens = append(smbTokens, token)
			smbMapper[token] = symbol
		}
	}
	currencyTokens := make([]string, 0, len(currencies))
	currencyMapper := make(map[string]string)
	for _, currency := range currencies {
		token, found := cg.currencyMapper[currency]
		if found {
			currencyTokens = append(currencyTokens, token)
			currencyMapper[token] = currency
		}
	}
	httpClient := http.DefaultClient
	return func(ctx context.Context) ([]price.Value, error) {
		records, err := coingecko.GetCurrentPrices(ctx, httpClient, smbTokens, currencyTokens)
		if err != nil {
			return nil, err
		}
		values := make([]price.Value, 0, len(records))
		for _, record := range records {
			values = append(values, price.Value{
				TimeStamp: time.Now(),
				Source:    "coingecko",
				Symbol:    smbMapper[record.Id],
				Currency:  currencyMapper[record.Currency],
				Price:     record.Value,
			})
		}
		return values, nil
	}
}
