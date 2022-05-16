package prices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type CoinGecko struct {
	symbols    map[string]string `yaml:"symbols"`
	currencies map[string]bool   `yaml:"currencies"`
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
func (cg *CoinGecko) GetterFunc(symbols map[string]string, currencies map[string]string) func(ctx context.Context) (map[string]map[string]float64, error) {
	requestSymbols := make([]string, 0, len(symbols))
	for symbol, _ := range symbols {
		requestSymbols = append(requestSymbols, symbol)
	}
	requestCurrencies := make([]string, 0, len(currencies))
	for currency, _ := range currencies {
		requestCurrencies = append(requestCurrencies, strings.ToLower(currency))
	}

	return getFunc(strings.Join(requestSymbols, ","), strings.Join(requestCurrencies, ","), symbols, currencies)
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
