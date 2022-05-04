package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
)

type PricesCache struct {
	cache  cache.Wrapper
	next   SymbolPrices
	expiry time.Duration
}

func NewPricesCache(cache cache.Wrapper, expiry time.Duration) *PricesCache {
	return &PricesCache{cache: cache, expiry: expiry}
}
func (cp *PricesCache) Wrap(next SymbolPrices) *PricesCache {
	cp.next = next
	return cp
}
func (cp *PricesCache) SetSymbolPrices(ctx context.Context, symbol string, source string, timeStamp time.Time, prices map[string]float64) error {
	if cp.next != nil {
		if err := cp.next.SetSymbolPrices(ctx, symbol, source, timeStamp, prices); err != nil {
			return fmt.Errorf("%w: %s", ErrWriting, err)
		}
		if err := cacheUnsetSymbolPrices(ctx, cp.cache, symbol); err != nil {
			return fmt.Errorf("%w: %s", ErrWriting, err)
		}
		return nil
	}
	record, err := cacheGetSymbolPrices(ctx, cp.cache, symbol)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			return fmt.Errorf("%w: %s", ErrReading, err)
		}
		record = NewPricesRecord(symbol, source, timeStamp)
	}
	record.Prices = prices
	if err := cacheSetSymbolPrices(ctx, cp.cache, record, cp.expiry); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func (cp *PricesCache) GetSymbolPrices(ctx context.Context, symbol string) (*PricesRecord, error) {
	value, err := cacheGetSymbolPrices(ctx, cp.cache, symbol)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrReading, err)
		}
		if cp.next == nil {
			return nil, ErrNotFound
		}
		value, err = cp.next.GetSymbolPrices(ctx, symbol)
		if err != nil {
			return nil, err
		}
		if err := cacheSetSymbolPrices(ctx, cp.cache, value, cp.expiry); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrWriting, err)
		}
	}
	return value, nil
}

func cacheUnsetSymbolPrices(ctx context.Context, cache cache.Wrapper, symbol string) error {
	if err := cache.Unset(ctx, toPricesCacheKey(symbol)); err != nil {
		return err
	}
	return nil
}
func cacheSetSymbolPrices(ctx context.Context, cache cache.Wrapper, prices *PricesRecord, expiry time.Duration) error {
	if err := cache.Set(ctx, toPricesCacheKey(prices.Symbol), prices, expiry); err != nil {
		return err
	}
	return nil
}
func cacheGetSymbolPrices(ctx context.Context, cache cache.Wrapper, symbol string) (*PricesRecord, error) {
	var value PricesRecord
	if err := cache.Get(ctx, toPricesCacheKey(symbol), &value); err != nil {
		return nil, err
	}
	return &value, nil
}
func toPricesCacheKey(symbol string) string {
	return symbol
}