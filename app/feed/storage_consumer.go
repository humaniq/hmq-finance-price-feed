package feed

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
)

type StorageConsumer struct {
	back        storage.PricesSaver
	in          chan []*state.PriceValue
	next        []Consumer
	done        chan interface{}
	leasesMutex sync.Mutex
	leases      int
	state       map[string]*state.Prices
	name        string
}

func NewStorageConsumer(name string, backend storage.PricesSaver, pricesState map[string]*state.Prices) *StorageConsumer {
	consumer := &StorageConsumer{
		back:  backend,
		in:    make(chan []*state.PriceValue),
		done:  make(chan interface{}),
		state: pricesState,
		name:  name,
	}
	return consumer
}
func (sc *StorageConsumer) In() chan<- []*state.PriceValue {
	return sc.in
}
func (sc *StorageConsumer) WithNext(next ...Consumer) *StorageConsumer {
	sc.next = append(sc.next, next...)
	for _, n := range next {
		n.Lease()
	}
	return sc
}
func (sc *StorageConsumer) WaitForDone() {
	<-sc.done
}
func (sc *StorageConsumer) Lease() chan<- []*state.PriceValue {
	sc.leasesMutex.Lock()
	defer sc.leasesMutex.Unlock()
	sc.leases++
	return sc.In()
}
func (sc *StorageConsumer) Release() {
	sc.leasesMutex.Lock()
	defer sc.leasesMutex.Unlock()
	sc.leases--
	if sc.leases == 0 {
		close(sc.in)
	}
}

func (sc *StorageConsumer) Run() {
	defer close(sc.done)
	defer func() {
		for _, next := range sc.next {
			next.Release()
		}
	}()
	ctx := context.Background()
	for prices := range sc.in {
		for _, price := range prices {
			currencyPrices, found := sc.state[price.Currency]
			if !found {
				logger.Info(ctx, "[%s] no currency available, skipping: %+v", sc.name, price)
				continue
			}
			currencyPrices.Commit(price)
		}
		var nextItems []*state.PriceValue
		for currency, currencyPrices := range sc.state {
			if len(currencyPrices.Changes()) > 0 {
				if err := sc.back.SavePrices(ctx, currency, currencyPrices); err != nil {
					logger.Error(ctx, "[%s] error saving prices: %s", sc.name, err)
					continue
				}
				nextItems = append(nextItems, currencyPrices.Changes()...)
				for _, priceChange := range currencyPrices.Changes() {
					logger.Info(ctx, "[%s] price changed: %+v", sc.name, priceChange)
				}
				currencyPrices.Stage()
			}
		}
		if len(nextItems) > 0 {
			for _, next := range sc.next {
				next.In() <- nextItems
			}
		}
	}
}
