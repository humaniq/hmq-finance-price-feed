package feed

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
)

type Provider interface {
	Provide(ctx context.Context, feed chan<- []*state.PriceValue) error
}

type Releaser interface {
	Lease() chan<- []*state.PriceValue
	Release()
}

type Waiter interface {
	WaitForDone()
}

type Consumer interface {
	In() chan<- []*state.PriceValue
	Waiter
	Releaser
}
