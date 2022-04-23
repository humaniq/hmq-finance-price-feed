package svc

import (
	"context"
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type FeederFunc func(record PriceRecord) error

type PriceFeed struct {
	prices    chan PriceRecord
	done      chan error
	consumers map[string]chan PriceRecord
}

func NewPriceFeed() *PriceFeed {
	return &PriceFeed{
		prices:    make(chan PriceRecord),
		consumers: make(map[string]chan PriceRecord),
	}
}
func (pf *PriceFeed) Queue() chan PriceRecord {
	return pf.prices
}
func (pf *PriceFeed) WithConsumerChan(name string, consumer chan PriceRecord) *PriceFeed {
	pf.consumers[name] = consumer
	return pf
}
func (pf *PriceFeed) Start() error {
	go pf.run()
	return nil
}
func (pf *PriceFeed) Stop() error {
	close(pf.prices)
	return nil
}
func (pf *PriceFeed) WaitForDone() error {
	return <-pf.done
}
func (pf *PriceFeed) run() {
	ctx := context.Background()
	defer close(pf.done)
	for record := range pf.prices {
		var wg sync.WaitGroup
		wg.Add(len(pf.consumers))
		for name, queue := range pf.consumers {
			logger.Info(ctx, "ENQUEUING to %s: %+v", name, record)
			go func(record PriceRecord, queue chan PriceRecord) {
				defer wg.Done()
				queue <- record
			}(record, queue)
		}
		wg.Wait()
	}
}
