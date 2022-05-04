package prices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"net/http"
	"os"
	"strings"
)

type CoinGecko struct {
	symbols    map[string]string
	currencies map[string]bool
}

func CoinGeckoFromFile(configFilePath string) (*CoinGecko, error) {
	symbolsFile, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer symbolsFile.Close()

	var config struct {
		Symbols    map[string]string
		Currencies []string
	}
	if err := json.NewDecoder(symbolsFile).Decode(&config); err != nil {
		return nil, err
	}
	currencies := make(map[string]bool)
	for _, key := range config.Currencies {
		currencies[key] = true
	}
	return &CoinGecko{
		symbols:    config.Symbols,
		currencies: currencies,
	}, nil
}
func (cg *CoinGecko) GetterFunc(symbols []string, currencies []string) func(ctx context.Context) (map[string]map[string]float64, error) {
	ctx := context.Background()
	requestSymbols := make(map[string]string)
	requestIds := make([]string, 0, len(symbols))
	for _, symbol := range symbols {
		id, found := cg.symbols[strings.ToLower(symbol)]
		if !found {
			logger.Error(ctx, "symbol not found: %s", symbol)
			continue
		}
		requestSymbols[id] = symbol
		requestIds = append(requestIds, id)
	}
	currenciesMapper := make(map[string]string)
	requestCurrencies := make([]string, 0, len(currencies))
	for _, currency := range currencies {
		if cg.currencies[strings.ToLower(currency)] {
			requestCurrencies = append(requestCurrencies, strings.ToLower(currency))
			currenciesMapper[strings.ToLower(currency)] = currency
		} else {
			logger.Error(ctx, "currency not found: %s", currency)
		}
	}

	return getFunc(strings.Join(requestIds, ","), strings.Join(requestCurrencies, ","), requestSymbols, currenciesMapper)
}

func getFunc(ids string, vsCurrencies string, symbolsMapper map[string]string, currenciesMapper map[string]string) func(ctx context.Context) (map[string]map[string]float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s", ids, vsCurrencies)
	return func(ctx context.Context) (map[string]map[string]float64, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("wrong response status")
		}
		defer resp.Body.Close()
		resultData := make(map[string]map[string]float64)
		if err := json.NewDecoder(resp.Body).Decode(&resultData); err != nil {
			return nil, err
		}
		result := make(map[string]map[string]float64)
		for key, val := range resultData {
			mappedRecord := make(map[string]float64)
			for currency, value := range val {
				mappedRecord[currenciesMapper[currency]] = value
			}
			result[symbolsMapper[key]] = mappedRecord
		}
		return result, nil
	}
}
