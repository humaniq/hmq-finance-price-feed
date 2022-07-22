package svc

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

var ErrNoValue = errors.New("no value")

type PriceRecord struct {
	Source        string    `json:"source"`
	Symbol        string    `json:"symbol"`
	Currency      string    `json:"currency"`
	Price         float64   `json:"price"`
	PreviousPrice float64   `json:"previousPrice,omitempty"`
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
	ds    *gds.Client
}

func NewPriceSvc() *PriceSvc {
	return &PriceSvc{}
}
func (ps *PriceSvc) WithCache(cache cache.Wrapper) *PriceSvc {
	ps.cache = cache
	return ps
}
func (ps *PriceSvc) WithGDSClient(ds *gds.Client) *PriceSvc {
	ps.ds = ds
	return ps
}
func (ps *PriceSvc) SetSymbolPrice(ctx context.Context, symbol string, currency string, price float64, source string) error {
	current, err := ps.GetLatestSymbolPrice(ctx, symbol, currency)
	if err != nil && !errors.Is(err, ErrNoValue) {
		return err
	}
	priceValue := NewPriceRecord(symbol, currency, price, source)
	if current != nil {
		priceValue = priceValue.WithPreviousPrice(current.Price)
	}
	if ps.ds != nil {
		if err := ps.ds.Write(ctx, toPriceKey(symbol, currency), priceValue); err != nil {
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
func (ps *PriceSvc) GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error) {
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
