package storage

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
)

type PricesSaver interface {
	SavePrices(ctx context.Context, key string, value *state.AssetPrices) error
}
type PricesLoader interface {
	LoadPrices(ctx context.Context, key string) (*state.AssetPrices, error)
}
type Prices interface {
	PricesLoader
	PricesSaver
}
