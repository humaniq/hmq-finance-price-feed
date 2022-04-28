package svc

import (
	"context"
	"errors"
	"time"
)

type PriceStateEstimateWrapper struct {
	backend Pricer
	ethalon string
}

func NewPriceStateEstimateWrapper(backend Pricer, ethalon string) *PriceStateEstimateWrapper {
	return &PriceStateEstimateWrapper{backend: backend, ethalon: ethalon}
}
func (psew *PriceStateEstimateWrapper) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
	value, err := psew.backend.GetLatestSymbolPrice(ctx, symbol, currency)
	if err != nil {
		if !errors.Is(err, ErrNoValue) {
			return nil, err
		}
		symbolPrice, err := psew.backend.GetLatestSymbolPrice(ctx, symbol, psew.ethalon)
		if err != nil {
			return nil, err
		}
		currencyPrice, err := psew.backend.GetLatestSymbolPrice(ctx, currency, psew.ethalon)
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

func (psew *PriceStateEstimateWrapper) SetSymbolPrice(ctx context.Context, price *PriceRecord) error {
	return psew.backend.SetSymbolPrice(ctx, price)
}
