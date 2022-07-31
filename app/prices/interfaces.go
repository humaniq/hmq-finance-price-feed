package prices

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type ProviderFunc func(ctx context.Context) ([]price.Value, error)
