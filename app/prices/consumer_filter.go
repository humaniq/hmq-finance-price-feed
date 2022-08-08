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

type ConsumerWorkerFilterWrapper struct {
	filterFunc ConsumerFilterFunc
	worker     ConsumerWorker
}

func NewConsumerWorkerFilterWpapper(filter ConsumerFilterFunc) *ConsumerWorkerFilterWrapper {
	return &ConsumerWorkerFilterWrapper{filterFunc: filter}
}
func (cwfw *ConsumerWorkerFilterWrapper) Wrap(worker ConsumerWorker) *ConsumerWorkerFilterWrapper {
	cwfw.worker = worker
	return cwfw
}
func (cwfw *ConsumerWorkerFilterWrapper) Work(ctx context.Context, values []price.Value) error {
	var filtered []price.Value
	for _, value := range values {
		if cwfw.filterFunc(ctx, value) {
			filtered = append(filtered, value)
		}
	}
	if len(filtered) == 0 {
		return nil
	}
	return cwfw.worker.Work(ctx, filtered)
}
