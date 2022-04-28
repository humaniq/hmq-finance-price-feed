package svc

import (
	"context"
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type StoreConsumer struct {
	back        Pricer
	in          chan *FeedItem
	next        []AsyncConsumer
	done        chan interface{}
	filters     []PriceFilterFunc
	leasesMutex sync.Mutex
	leases      int
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
func (sh *StoreConsumer) WithNext(next ...AsyncConsumer) *StoreConsumer {
	sh.next = append(sh.next, next...)
	for _, n := range next {
		n.Lease()
	}
	return sh
}
func (sh *StoreConsumer) WithFilters(fn ...PriceFilterFunc) *StoreConsumer {
	sh.filters = append(sh.filters, fn...)
	return sh
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
func (sh *StoreConsumer) Lease() chan<- *FeedItem {
	sh.leasesMutex.Lock()
	defer sh.leasesMutex.Unlock()
	sh.leases++
	return sh.In()
}
func (sh *StoreConsumer) Release() {
	sh.leasesMutex.Lock()
	defer sh.leasesMutex.Unlock()
	sh.leases--
	if sh.leases == 0 {
		sh.Stop()
	}
}

func (sh *StoreConsumer) run(ctx context.Context, in <-chan *FeedItem) error {
	defer close(sh.done)
	defer func() {
		for _, next := range sh.next {
			next.Release()
		}
	}()
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
		for _, next := range sh.next {
			next.In() <- &FeedItem{records: records}
		}
	}
	return nil
}
