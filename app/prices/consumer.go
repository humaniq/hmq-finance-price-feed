package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"sync"
)

type ConsumerWorker interface {
	Work(ctx context.Context, values []price.Value)
}

type Consumer struct {
	done    chan interface{}
	workers []ConsumerWorker
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

func (c *Consumer) Run(ctx context.Context, in <-chan []price.Value) {
	defer close(c.done)
	for values := range in {
		var wg sync.WaitGroup
		wg.Add(len(c.workers))
		for _, worker := range c.workers {
			go func(w ConsumerWorker) {
				defer wg.Done()
				w.Work(ctx, values)
			}(worker)
		}
		wg.Wait()
	}
}
