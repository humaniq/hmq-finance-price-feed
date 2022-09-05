package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"sync"
)

type ConsumerWorker interface {
	Work(ctx context.Context, values []price.Value) error
}

type Consumer struct {
	done     chan interface{}
	workers  []ConsumerWorker
	enriches []ConsumerEnrichFunc
	filters  []ConsumerFilterFunc
}

func NewConsumer() *Consumer {
	return &Consumer{done: make(chan interface{})}
}
func (c *Consumer) AddWorker(w ConsumerWorker) {
	c.workers = append(c.workers, w)
}

func (c *Consumer) Consume(ctx context.Context, in <-chan []price.Value) error {
	go c.Run(ctx, in)
	return nil
}
func (c *Consumer) WaitForDone() {
	<-c.done
}

func (c *Consumer) WithEnrich(fn ...ConsumerEnrichFunc) *Consumer {
	c.enriches = append(c.enriches, fn...)
	return c
}
func (c *Consumer) WithFilters(fn ...ConsumerFilterFunc) *Consumer {
	c.filters = append(c.filters, fn...)
	return c
}

func (c *Consumer) Run(ctx context.Context, in <-chan []price.Value) {
	defer close(c.done)
	for values := range in {
		workValues := make([]price.Value, 0, len(values))
		for _, enrich := range c.enriches {
			for _, value := range values {
				workValues = append(workValues, enrich(ctx, value)...)
			}
		}
		var filteredValues []price.Value
		for _, filter := range c.filters {
			for _, value := range workValues {
				if filter(ctx, value) {
					filteredValues = append(filteredValues, value)
				}
			}
		}
		var wg sync.WaitGroup
		wg.Add(len(c.workers))
		for _, worker := range c.workers {
			go func(w ConsumerWorker) {
				defer wg.Done()
				if err := w.Work(ctx, filteredValues); err != nil {
					app.Logger().Error(ctx, "WORKER ERROR: %s", err)
				}
			}(worker)
		}
		wg.Wait()
	}
}
