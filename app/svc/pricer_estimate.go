package svc

import (
	"context"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"time"
)

type PriceStorageEstimateWrapper struct {
	backend  storage.SymbolPrices
	standard string
}

func NewPriceStorageEstimateWrapper(backend storage.SymbolPrices, standard string) *PriceStorageEstimateWrapper {
	return &PriceStorageEstimateWrapper{backend: backend, standard: standard}
}

func (psew *PriceStorageEstimateWrapper) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
	value, err := psew.backend.GetLatestSymbolPrice(ctx, symbol, currency)
	if err != nil {
		if !errors.Is(err, ErrNoValue) {
			return nil, err
		}
		symbolPrice, err := psew.backend.GetLatestSymbolPrice(ctx, symbol, psew.standard)
		if err != nil {
			return nil, err
		}
		currencyPrice, err := psew.backend.GetLatestSymbolPrice(ctx, currency, psew.standard)
		if err != nil {
			return nil, err
		}
		return &PriceRecord{
			Source:        "estimation",
			Symbol:        symbol,
			Currency:      currency,
			Price:         symbolPrice.Price / currencyPrice.Price,
			PreviousPrice: 0,
			TimeStamp:     time.Now(),
		}, nil

	}
	return value, nil
}

func (psew *PriceStorageEstimateWrapper) SetSymbolPrice(ctx context.Context, price *PriceRecord) error {
	return psew.backend.SetSymbolPrice(ctx, price)
}
