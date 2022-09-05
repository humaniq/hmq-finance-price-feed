package prices

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type ConsumerEnrichFunc func(ctx context.Context, value price.Value) []price.Value

type ConsumerEnricherWrapper struct {
	worker ConsumerWorker
}

func NewConsumerEnricherWrapper(mapperCurrency string, assetCurrencies []string) *ConsumerEnricherWrapper {
	return &ConsumerEnricherWrapper{}
}
func (cew *ConsumerEnricherWrapper) Wrap(worker ConsumerWorker) *ConsumerEnricherWrapper {
	cew.worker = worker
	return cew
}
