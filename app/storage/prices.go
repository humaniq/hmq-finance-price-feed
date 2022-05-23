package storage

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type PricesSaver interface {
	SavePrices(ctx context.Context, key string, value *price.Asset) error
}
type PricesLoader interface {
	LoadPrices(ctx context.Context, key string) (*price.Asset, error)
}
type Prices interface {
	PricesLoader
	PricesSaver
}
