package storage

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
)

type PricesSaver interface {
	SavePrices(ctx context.Context, key string, value *state.Prices) error
}
type PricesLoader interface {
	LoadPrices(ctx context.Context, key string) (*state.Prices, error)
}
type Prices interface {
	PricesLoader
	PricesSaver
}
