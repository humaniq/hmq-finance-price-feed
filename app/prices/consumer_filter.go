package prices

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type ConsumerFilterFunc func(ctx context.Context, value price.Value) bool

func AnyOf(fns ...ConsumerFilterFunc) ConsumerFilterFunc {
	return func(ctx context.Context, value price.Value) bool {
		for _, fn := range fns {
			if fn(ctx, value) {
				return true
			}
		}
		return false
	}
}
func AllOf(fns ...ConsumerFilterFunc) ConsumerFilterFunc {
	return func(ctx context.Context, value price.Value) bool {
		for _, fn := range fns {
			if !fn(ctx, value) {
				return false
			}
		}
		return true
	}
}
