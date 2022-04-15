package svc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
)

var ErrNoValue = errors.New("no value")

type PriceRecord struct {
	Source        string    `json:"source"`
	Symbol        string    `json:"symbol"`
	Currency      string    `json:"currency"`
	Price         float64   `json:"price"`
	PreviousPrice float64   `json:"previousPrice"`
	TimeStamp     time.Time `json:"timeStamp"`
}

func NewPriceRecord(symbol string, currency string, price float64, source string) *PriceRecord {
	return &PriceRecord{
		Source:    source,
		Symbol:    symbol,
		Currency:  currency,
		Price:     price,
		TimeStamp: time.Now(),
	}
}
func (pr *PriceRecord) WithPreviousPrice(price float64) *PriceRecord {
	pr.PreviousPrice = price
	return pr
}

type PriceSvc struct {
	cache cache.Wrapper
}

func NewPriceSvc() *PriceSvc {
	return &PriceSvc{}
}
func (ps *PriceSvc) WithCache(cache cache.Wrapper) *PriceSvc {
	ps.cache = cache
	return ps
}
func (ps *PriceSvc) SetSymbolPrice(ctx context.Context, symbol string, currency string, price float64, source string) error {
	current := &PriceRecord{}
	key := toPriceCacheKey(symbol, currency)
	if err := ps.cache.Get(ctx, key, current); err != nil && !errors.Is(err, cache.ErrNotFound) {
		return err
	}
	if err := ps.cache.Set(ctx, key, NewPriceRecord(symbol, currency, price, source).WithPreviousPrice(current.Price), time.Minute); err != nil {
		return err
	}
	return nil
}
func (ps *PriceSvc) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
	if ps.cache != nil {
		record, err := cacheGetSymbolPrice(ctx, ps.cache, symbol, currency)
		if err != nil && !errors.Is(err, cache.ErrNotFound) {
			return nil, err
		}
		if err == nil {
			return record, nil
		}
	}

	record := &PriceRecord{}
	if err := ps.cache.Get(ctx, toPriceCacheKey(symbol, currency), record); err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return nil, ErrNoValue
		}
		return nil, err
	}
	return record, nil
}

func toPriceCacheKey(symbol string, currency string) string {
	return fmt.Sprintf("%s-%s", symbol, currency)
}

func cacheSetSymbolPrice(ctx context.Context, cache cache.Wrapper, record *PriceRecord) error {
	return cache.Set(ctx, toPriceCacheKey(record.Symbol, record.Currency), record, time.Minute)
}
func cacheGetSymbolPrice(ctx context.Context, cache cache.Wrapper, symbol string, currency string) (*PriceRecord, error) {
	record := &PriceRecord{}
	if err := cache.Get(ctx, toPriceCacheKey(symbol, currency), record); err != nil {
		return nil, err
	}
	return record, nil
}
