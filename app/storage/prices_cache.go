package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
)

type PricesCache struct {
	cache  cache.Wrapper
	next   Prices
	expiry time.Duration
}

func NewPricesCache(cache cache.Wrapper, expiry time.Duration) *PricesCache {
	return &PricesCache{cache: cache, expiry: expiry}
}
func (cp *PricesCache) Wrap(next Prices) *PricesCache {
	cp.next = next
	return cp
}

func (cp *PricesCache) SavePrices(ctx context.Context, key string, value *state.AssetPrices) error {
	if cp.next != nil {
		if err := cp.next.SavePrices(ctx, key, value); err != nil {
			return fmt.Errorf("%w: %s", ErrWriting, err)
		}
	}
	if err := cacheSetPrices(ctx, cp.cache, key, value, cp.expiry); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}

func (cp *PricesCache) LoadPrices(ctx context.Context, key string) (*state.AssetPrices, error) {
	value, err := cacheGetPrices(ctx, cp.cache, key)
	if err != nil {
		if !errors.Is(err, cache.ErrNotFound) {
			return nil, fmt.Errorf("%w: %s", ErrReading, err)
		}
		if cp.next == nil {
			return nil, ErrNotFound
		}
		value, err = cp.next.LoadPrices(ctx, key)
		if err != nil {
			return nil, err
		}
		if err := cacheSetPrices(ctx, cp.cache, key, value, cp.expiry); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrWriting, err)
		}
	}
	return value, nil
}

func cacheUnsetPrices(ctx context.Context, cache cache.Wrapper, key string) error {
	if err := cache.Unset(ctx, toPricesCacheKey(key)); err != nil {
		return err
	}
	return nil
}
func cacheSetPrices(ctx context.Context, cache cache.Wrapper, key string, prices *state.AssetPrices, expiry time.Duration) error {
	if err := cache.Set(ctx, toPricesCacheKey(key), prices, expiry); err != nil {
		return err
	}
	return nil
}
func cacheGetPrices(ctx context.Context, cache cache.Wrapper, key string) (*state.AssetPrices, error) {
	value := state.NewPrices(key)
	if err := cache.Get(ctx, toPricesCacheKey(key), value); err != nil {
		return nil, err
	}
	return value, nil
}
func toPricesCacheKey(key string) string {
	return key
}
