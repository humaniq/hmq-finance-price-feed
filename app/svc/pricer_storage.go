package svc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

type PriceStateSvc struct {
	cache cache.Wrapper
	ds    *gds.Client
}

func NewPriceStateSvc() *PriceStateSvc {
	return &PriceStateSvc{}
}
func (ps *PriceStateSvc) WithCache(cache cache.Wrapper) *PriceStateSvc {
	ps.cache = cache
	return ps
}
func (ps *PriceStateSvc) WithGDSClient(ds *gds.Client) *PriceStateSvc {
	ps.ds = ds
	return ps
}
func (ps *PriceStateSvc) SetSymbolPrice(ctx context.Context, price *PriceRecord) error {
	current, err := ps.GetLatestSymbolPrice(ctx, price.Symbol, price.Currency)
	if err != nil && !errors.Is(err, ErrNoValue) {
		return err
	}
	priceValue := NewPriceRecord(price.Symbol, price.Currency, price.Price, price.Source)
	if current != nil {
		priceValue = priceValue.WithPreviousPrice(current.Price)
	}
	if ps.ds != nil {
		if err := ps.ds.Write(ctx, toPriceKey(price.Symbol, price.Currency), priceValue); err != nil {
			return err
		}
	}
	if ps.cache != nil {
		if err := cacheSetSymbolPrice(ctx, ps.cache, priceValue); err != nil {
			return err
		}
	}
	return nil
}
func (ps *PriceStateSvc) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
	needsCacheRefresh := false
	if ps.cache != nil {
		record, err := cacheGetSymbolPrice(ctx, ps.cache, symbol, currency)
		if err != nil {
			if !errors.Is(err, cache.ErrNotFound) {
				return nil, err
			} else {
				needsCacheRefresh = true
			}
		}
		if err == nil {
			return record, nil
		}
	}
	if ps.ds != nil {
		record := &PriceRecord{}
		if err := ps.ds.Read(ctx, toPriceKey(symbol, currency), record); err != nil {
			if errors.Is(err, gds.ErrNotFound) {
				return nil, ErrNoValue
			}
			return nil, err
		}
		if needsCacheRefresh {
			if err := cacheSetSymbolPrice(ctx, ps.cache, record); err != nil {
				return nil, err
			}
		}
		return record, nil
	}

	return nil, ErrNoValue
}

func toPriceKey(symbol string, currency string) string {
	return fmt.Sprintf("%s-%s", symbol, currency)
}

func cacheSetSymbolPrice(ctx context.Context, cache cache.Wrapper, record *PriceRecord) error {
	return cache.Set(ctx, toPriceKey(record.Symbol, record.Currency), record, time.Minute)
}
func cacheGetSymbolPrice(ctx context.Context, cache cache.Wrapper, symbol string, currency string) (*PriceRecord, error) {
	record := &PriceRecord{}
	if err := cache.Get(ctx, toPriceKey(symbol, currency), record); err != nil {
		return nil, err
	}
	return record, nil
}
