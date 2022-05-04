package svc

import (
	"context"
	"errors"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"time"
)

type PricesStorage struct {
	backend  storage.SymbolPricesGetter
	required []string
	std      string
}

func NewPricesStorage(backend storage.SymbolPricesGetter) *PricesStorage {
	return &PricesStorage{backend: backend}
}
func (ps *PricesStorage) WithRequiredCurrencies(required ...string) *PricesStorage {
	for _, item := range required {
		ps.required = append(ps.required, item)
	}
	return ps
}
func (ps *PricesStorage) WithStdCurrency(std string) *PricesStorage {
	ps.std = std
	return ps
}

func (ps *PricesStorage) GetPrices(ctx context.Context, symbols []string, currencies []string) (map[string]SymbolPrices, error) {
	result := make(map[string]SymbolPrices)
	missing := newMissingPrices()
	for _, symbol := range symbols {
		symbolPrices, err := ps.backend.GetSymbolPrices(ctx, symbol)
		if err != nil {
			if !errors.Is(err, storage.ErrNotFound) {
				return nil, err
			}
			continue
		}
		prices := make(map[string]SymbolPrice)
		for _, currency := range currencies {
			price, found := symbolPrices.Prices[currency]
			if !found {
				if ps.std != "" && symbolPrices.Prices[ps.std] != 0 {
					missing.addPrice(missingPrice{
						Symbol:        symbol,
						Currency:      currency,
						StandardPrice: symbolPrices.Prices[ps.std],
					})
				}
				continue
			}
			prices[currency] = SymbolPrice{
				Source:    symbolPrices.Source,
				Value:     price,
				TimeStamp: symbolPrices.TimeStamp,
			}
		}
		result[symbol] = prices
	}
	if !missing.isEmpty() {
		stdPricesData, err := ps.backend.GetSymbolPrices(ctx, fmt.Sprintf("std_%s", ps.std))
		if err != nil && !errors.Is(err, storage.ErrNotFound) {
			return nil, err
		}
		stdPrices := stdPricesData.Prices
		for _, miss := range missing.missing {
			stdPrice, found := stdPrices[miss.Currency]
			if !found {
				continue
			}
			result[miss.Symbol][miss.Currency] = SymbolPrice{
				Source:    "estimation",
				Value:     miss.StandardPrice / stdPrice,
				TimeStamp: time.Now(),
			}
		}
	}
	return result, nil
}

type missingPrices struct {
	missing    []missingPrice
	currencies map[string]bool
}

func newMissingPrices() *missingPrices {
	return &missingPrices{
		currencies: make(map[string]bool),
	}
}
func (mp *missingPrices) addPrice(prices ...missingPrice) {
	for _, price := range prices {
		mp.missing = append(mp.missing, price)
		mp.currencies[price.Currency] = true
	}
}
func (mp *missingPrices) isEmpty() bool {
	return len(mp.missing) == 0
}
func (mp *missingPrices) Currencies() []string {
	value := make([]string, 0, len(mp.currencies))
	for key, _ := range mp.currencies {
		value = append(value, key)
	}
	return value
}

type missingPrice struct {
	Symbol        string
	Currency      string
	StandardPrice float64
}
