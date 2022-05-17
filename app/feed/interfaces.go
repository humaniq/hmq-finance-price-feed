package feed

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type Provider interface {
	Provide(ctx context.Context, feed chan<- []price.Value) error
}

type Releaser interface {
	Lease() chan<- []price.Value
	Release()
}

type Waiter interface {
	WaitForDone()
}

type Consumer interface {
	In() chan<- []price.Value
	Waiter
	Releaser
}
