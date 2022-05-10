package feed

import (
	"context"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type StoreConsumer struct {
	back             storage.Prices
	in               chan []*state.Price
	next             []Consumer
	done             chan interface{}
	leasesMutex      sync.Mutex
	leases           int
	symbolsFinder    map[string]bool
	currenciesFinder map[string]bool
}

func NewStoreConsumer(backend storage.Prices) *StoreConsumer {
	consumer := &StoreConsumer{
		back: backend,
		in:   make(chan []*state.Price),
		done: make(chan interface{}),
	}
	go consumer.run(context.Background(), consumer.in)
	return consumer
}
func (sh *StoreConsumer) WithSymbols(symbols ...string) *StoreConsumer {
	if sh.symbolsFinder == nil {
		sh.symbolsFinder = make(map[string]bool)
	}
	for _, s := range symbols {
		sh.symbolsFinder[s] = true
	}
	return sh
}
func (sh *StoreConsumer) WithCurrencies(currencies ...string) *StoreConsumer {
	if sh.currenciesFinder == nil {
		sh.currenciesFinder = make(map[string]bool)
	}
	for _, c := range currencies {
		sh.currenciesFinder[c] = true
	}
	return sh
}
func (sh *StoreConsumer) currencies() map[string]bool {
	if sh.currenciesFinder != nil {
		return sh.currenciesFinder
	}
	return config.CurrenciesKnown
}
func (sh *StoreConsumer) symbols() map[string]bool {
	if sh.symbolsFinder != nil {
		return sh.symbolsFinder
	}
	smb := make(map[string]bool)
	for s, _ := range config.SymbolsKnown {
		smb[s] = true
	}
	return smb
}
func (sh *StoreConsumer) In() chan<- []*state.Price {
	return sh.in
}
func (sh *StoreConsumer) WithNext(next ...Consumer) *StoreConsumer {
	sh.next = append(sh.next, next...)
	for _, n := range next {
		n.Lease()
	}
	return sh
}
func (sh *StoreConsumer) WaitForDone() {
	<-sh.done
}
func (sh *StoreConsumer) Lease() chan<- []*state.Price {
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
		close(sh.in)
	}
}

func (sh *StoreConsumer) run(ctx context.Context, in <-chan []*state.Price) {
	defer close(sh.done)
	defer func() {
		for _, next := range sh.next {
			next.Release()
		}
	}()
	symbols := sh.symbols()
	currencies := sh.currencies()
	mapper := make(map[string]*state.Prices)
	for currency, _ := range currencies {
		prices, err := sh.back.LoadPrices(ctx, currency)
		if err != nil {
			if !errors.Is(err, storage.ErrNotFound) {
				logger.Fatal(ctx, "StoreConsumer: error on prices loading: %s", err.Error())
				return
			}
			prices = state.NewPrices(currency).WithCommitFilters(
				state.CommitTimestampFilterFunc(),
				state.CommitSymbolsFilterFunc(symbols),
				state.CommitCurrenciesFilterFunc(map[string]bool{currency: true}),
			)
		}
		mapper[currency] = prices
	}
	for items := range in {
		if len(items) == 0 {
			logger.Warn(ctx, "StoreConsumer: items are empty")
			continue
		}
		for _, record := range items {
			if !currencies[record.Currency] {
				continue
			}
			prices, found := mapper[record.Currency]
			if !found {
				prices = state.NewPrices(record.Currency).WithCommitFilters(
					state.CommitTimestampFilterFunc(),
					state.CommitSymbolsFilterFunc(symbols),
					state.CommitCurrenciesFilterFunc(map[string]bool{record.Currency: true}),
				)
			}
			prices.Commit(record)
			mapper[record.Currency] = prices
		}
		var nexts []*state.Price
		for currency, prices := range mapper {
			changes := prices.Changes()
			if len(changes) > 0 {
				if err := sh.back.SavePrices(ctx, currency, prices); err != nil {
					logger.Critical(ctx, "StoreConsumer: error on save: %s", err.Error())
					continue
				}
				prices.Stage()
				nexts = append(nexts, changes...)
			}
		}
		if sh.next != nil && len(sh.next) > 0 && len(nexts) > 0 {
			for _, next := range sh.next {
				next.In() <- nexts
			}
		}
	}
}
