package svc

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type StoreConsumer struct {
	back    Pricer
	in      chan *FeedItem
	out     []chan<- *FeedItem
	done    chan interface{}
	filters []PriceFilterFunc
}

func NewStoreConsumer(backend Pricer) *StoreConsumer {
	return &StoreConsumer{
		back: backend,
		in:   make(chan *FeedItem),
		done: make(chan interface{}),
	}
}
func (sh *StoreConsumer) In() chan<- *FeedItem {
	return sh.in
}
func (sh *StoreConsumer) WithNext(next ...chan<- *FeedItem) *StoreConsumer {
	sh.out = append(sh.out, next...)
	return sh
}
func (sh *StoreConsumer) WithFilters(fn ...PriceFilterFunc) {
	sh.filters = append(sh.filters, fn...)
}

func (sh *StoreConsumer) Start() error {
	go sh.run(context.Background(), sh.in)
	return nil
}
func (sh *StoreConsumer) Stop() error {
	close(sh.in)
	return nil
}
func (sh *StoreConsumer) WaitForDone() {
	<-sh.done
}
func (sh *StoreConsumer) Consume(ctx context.Context, in <-chan *FeedItem) error {
	return sh.run(ctx, in)
}

func (sh *StoreConsumer) run(ctx context.Context, in <-chan *FeedItem) error {
	defer close(sh.done)
	for item := range in {
		if len(item.records) == 0 {
			logger.Warn(ctx, "items are empty")
			continue
		}
		records := make([]*PriceRecord, 0, len(item.records))
		for _, filter := range sh.filters {
			item.Filter(ctx, filter)
		}
		for _, record := range item.records {
			if err := sh.back.SetSymbolPrice(ctx, record); err != nil {
				logger.Error(ctx, "error setting symbol price: %s", err.Error())
				continue
			}
			records = append(records, record)
		}
		for _, next := range sh.out {
			next <- &FeedItem{records: records}
		}
	}
	return nil
}
