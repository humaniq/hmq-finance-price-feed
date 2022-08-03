package storage

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type AssetInMemoryRecord struct {
	asset   *price.Asset
	expires time.Time
}

type InMemory struct {
	mx     sync.RWMutex
	assets map[string]AssetInMemoryRecord
	next   Prices
	expiry time.Duration
}

func NewInMemory(expiry time.Duration) *InMemory {
	return &InMemory{
		assets: make(map[string]AssetInMemoryRecord),
		expiry: expiry,
	}
}
func (im *InMemory) Wrap(next Prices) *InMemory {
	im.next = next
	return im
}
func (im *InMemory) Warm(ctx context.Context, currencyList []string, rotationTimer time.Duration) *InMemory {
	if len(currencyList) == 0 {
		return im
	}
	go func() {
		ticker := time.NewTicker(rotationTimer / time.Duration(len(currencyList)))
		index := 0
		for range ticker.C {
			if im.next == nil {
				app.Logger().Error(ctx, "next is nil, no warm available")
				continue
			}
			currentAssetValue, err := im.next.LoadPrices(ctx, currencyList[index])
			if err != nil {
				app.Logger().Error(ctx, "WARM: failed to update %s: %s", currencyList[index], err)
				continue
			}
			im.set(currencyList[index], currentAssetValue)
			app.Logger().Info(ctx, "WARM: updated price for %s", currencyList[index])
			index++
			if index >= len(currencyList) {
				index = 0
			}
		}
	}()
	return im
}
func (im *InMemory) LoadPrices(ctx context.Context, key string) (*price.Asset, error) {
	record := im.get(key)
	if record.expires.Before(time.Now()) {
		if im.next != nil {
			value, err := im.next.LoadPrices(ctx, key)
			if err != nil {
				return nil, err
			}
			im.set(key, value)
			return value, nil
		}
		return nil, app.ErrNotFound
	}
	return record.asset, nil
}
func (im *InMemory) SavePrices(ctx context.Context, key string, value *price.Asset) error {
	if im.next != nil {
		if err := im.next.SavePrices(ctx, key, value); err != nil {
			return fmt.Errorf("%w: %s", ErrWriting, err)
		}
	}
	im.set(key, value)
	return nil
}

func (im *InMemory) set(key string, value *price.Asset) {
	asset := AssetInMemoryRecord{
		asset:   value,
		expires: time.Now().Add(im.expiry),
	}
	im.mx.Lock()
	defer im.mx.Unlock()
	im.assets[key] = asset
}
func (im *InMemory) get(key string) AssetInMemoryRecord {
	im.mx.RLock()
	defer im.mx.RUnlock()
	return im.assets[key]
}
