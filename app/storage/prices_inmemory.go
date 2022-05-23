package storage

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"sync"
	"time"

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
