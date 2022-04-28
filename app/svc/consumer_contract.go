package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"sync"
)

type ContractPricesConsumer struct {
	getter      PriceGetter
	setter      PriceSetter
	in          chan *FeedItem
	next        []AsyncConsumer
	done        chan interface{}
	filters     []PriceFilterFunc
	back        Pricer
	leasesMutex sync.Mutex
	leases      int
}

func NewContractPricesConsumer() *ContractPricesConsumer {
	return &ContractPricesConsumer{
		in:   make(chan *FeedItem),
		done: make(chan interface{}),
	}
}
func (cpc *ContractPricesConsumer) WithSetter(setter PriceSetter) *ContractPricesConsumer {
	cpc.setter = setter
	return cpc
}
func (cpc *ContractPricesConsumer) WithGetter(getter PriceGetter) *ContractPricesConsumer {
	cpc.getter = getter
	return cpc
}

func (cpc *ContractPricesConsumer) Lease() chan<- *FeedItem {
	cpc.leasesMutex.Lock()
	defer cpc.leasesMutex.Unlock()
	cpc.leases++
	return cpc.In()
}
func (cpc *ContractPricesConsumer) Release() {
	cpc.leasesMutex.Lock()
	defer cpc.leasesMutex.Unlock()
	cpc.leases--
	if cpc.leases == 0 {
		cpc.Stop()
	}
}

func (cpc *ContractPricesConsumer) In() chan<- *FeedItem {
	return cpc.in
}
func (cpc *ContractPricesConsumer) WithNext(next ...AsyncConsumer) *ContractPricesConsumer {
	cpc.next = append(cpc.next, next...)
	for _, n := range next {
		n.Lease()
	}
	return cpc
}
func (cpc *ContractPricesConsumer) WithFilters(fn ...PriceFilterFunc) *ContractPricesConsumer {
	cpc.filters = append(cpc.filters, fn...)
	return cpc
}

func (cpc *ContractPricesConsumer) Start() error {
	go cpc.run(context.Background(), cpc.in)
	return nil
}
func (cpc *ContractPricesConsumer) Stop() error {
	close(cpc.in)
	return nil
}
func (cpc *ContractPricesConsumer) WaitForDone() {
	<-cpc.done
}

func (cpc *ContractPricesConsumer) Consume(ctx context.Context, feed <-chan *FeedItem) error {
	return cpc.run(ctx, feed)
}

func (cpc *ContractPricesConsumer) run(ctx context.Context, feed <-chan *FeedItem) error {
	defer close(cpc.done)
	defer func() {
		for _, next := range cpc.next {
			next.Release()
		}
	}()
	for item := range feed {
		for _, filter := range cpc.filters {
			item.Filter(ctx, filter)
		}
		for _, price := range item.records {
			logger.Debug(ctx, "contract write call with %+v", price)
			if cpc.setter != nil {
				if err := cpc.setter.SetSymbolPrice(ctx, price); err != nil {
					logger.Error(ctx, "error setting contract price: %s", err)
					continue
				}
			}
		}
		for _, next := range cpc.next {
			next.In() <- item
		}
	}
	return nil
}
