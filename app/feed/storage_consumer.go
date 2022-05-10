package feed

import (
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
)

type StorageConsumer struct {
	back        storage.PricesSaver
	in          chan []*state.Price
	next        []Consumer
	done        chan interface{}
	leasesMutex sync.Mutex
	leases      int
	state       map[string]*state.Prices
}

func NewStorageConsumer(backend storage.PricesSaver, pricesState map[string]*state.Prices) *StorageConsumer {
	return &StorageConsumer{
		back:  backend,
		in:    make(chan []*state.Price),
		done:  make(chan interface{}),
		state: pricesState,
	}
}
func (sc *StorageConsumer) In() chan<- []*state.Price {
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
func (sc *StorageConsumer) Lease() chan<- []*state.Price {
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

func (sc *StorageConsumer) run() {

}
