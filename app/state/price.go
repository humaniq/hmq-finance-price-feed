package state

import (
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"sort"
	"time"
)

type Price struct {
	Current       *PriceValue
	History       []PriceHistory
	historyOffset time.Duration
}

func NewPrice() *Price {
	return &Price{}
}
func (p *Price) WithHistoryOffset(offset time.Duration) *Price {
	p.historyOffset = offset
	return p
}
func (p *Price) Set(value *PriceValue) {
	p.Current = value
	history := make([]PriceHistory, 0, len(p.History)+1)
	for _, hr := range p.History {
		if p.historyOffset != 0 && hr.TimeStamp.Before(value.TimeStamp.Add(-1*p.historyOffset)) {
			continue
		}
		history = append(history, PriceHistory{
			TimeStamp: hr.TimeStamp,
			Value:     hr.Value,
		})
	}
	history = append(history, PriceHistory{
		TimeStamp: value.TimeStamp,
		Value:     value.Price,
	})
	sort.Slice(history, func(i, j int) bool {
		return history[i].TimeStamp.Before(history[j].TimeStamp)
	})
	p.History = history
}

type PriceHistory struct {
	TimeStamp time.Time
	Value     float64
}
type PriceValue struct {
	TimeStamp time.Time
	Source    string
	Symbol    string
	Currency  string
	Price     float64
}

func NewPriceValue(source string, symbol string, currency string, price float64, timeStamp time.Time) *PriceValue {
	return &PriceValue{
		TimeStamp: timeStamp,
		Source:    source,
		Symbol:    symbol,
		Currency:  currency,
		Price:     price,
	}
}

type AssetPrices struct {
	Currency      string            `json:"currency"`
	Values        map[string]*Price `json:"values"`
	changed       []*PriceValue
	commitFilters []CommitFilterFunc
}

func NewPrices(currency string) *AssetPrices {
	return &AssetPrices{
		Currency: currency,
		Values:   make(map[string]*Price),
		changed:  []*PriceValue{},
	}
}
func (ps *AssetPrices) Key() string {
	return ps.Currency
}
func (ps *AssetPrices) WithCommitFilters(fn ...CommitFilterFunc) *AssetPrices {
	ps.commitFilters = append(ps.commitFilters, fn...)
	return ps
}
func (ps *AssetPrices) Commit(price *PriceValue) bool {
	current, found := ps.Values[price.Symbol]
	if !found {
		current = NewPrice()
	} else {
		for _, filter := range ps.commitFilters {
			if !filter(current.Current, price) {
				return false
			}
		}
	}
	current.Set(price)
	ps.Values[price.Symbol] = current
	ps.changed = append(ps.changed, price)
	return true
}
func (ps *AssetPrices) Stage() {
	ps.changed = []*PriceValue{}
}
func (ps *AssetPrices) Changes() []*PriceValue {
	return ps.changed
}
func (ps *AssetPrices) Get(symbol string, withHistory bool) *Price {
	symbolPrice, found := ps.Values[symbol]
	if !found {
		return nil
	}
	p := Price{
		Current: symbolPrice.Current,
	}
	if withHistory {
		p.History = symbolPrice.History
	}
	return &p
}
func (ps *AssetPrices) Estimate(symbol string, currency string, withHistory bool) *Price {
	symbolPrice, found := ps.Values[symbol]
	if !found {
		return nil
	}
	currencyPrice, found := ps.Values[currency]
	if !found {
		return nil
	}
	p := Price{
		Current: &PriceValue{
			TimeStamp: time.Now(),
			Source:    "estimate",
			Symbol:    symbol,
			Currency:  currency,
			Price:     symbolPrice.Current.Price / currencyPrice.Current.Price,
		},
	}
	if withHistory {
		p.History = estimateHistory(symbolPrice.History, currencyPrice.History)
	}
	return &p
}
func estimateHistory(symbolPrices []PriceHistory, currencyPrices []PriceHistory) []PriceHistory {
	type historyValue struct {
		isSymbol bool
		value    PriceHistory
	}
	history := make([]historyValue, 0, len(symbolPrices)+len(currencyPrices))
	for _, sp := range symbolPrices {
		history = append(history, historyValue{isSymbol: true, value: sp})
	}
	for _, cp := range currencyPrices {
		history = append(history, historyValue{value: cp})
	}
	sort.Slice(history, func(i, j int) bool {
		return history[i].value.TimeStamp.Before(history[j].value.TimeStamp)
	})
	result := make([]PriceHistory, 0, len(history))
	currentSymbolPrice := float64(0)
	currentCurrencyPrice := float64(0)
	for _, hValue := range history {
		if hValue.isSymbol && hValue.value.Value != 0 {
			currentSymbolPrice = hValue.value.Value
		}
		if !hValue.isSymbol && hValue.value.Value != 0 {
			currentCurrencyPrice = hValue.value.Value
		}
		if currentCurrencyPrice != 0 && currentSymbolPrice != 0 {
			result = append(result, PriceHistory{
				TimeStamp: hValue.value.TimeStamp,
				Value:     currentSymbolPrice / currentCurrencyPrice,
			})
		}
	}
	return result
}

type CommitFilterFunc func(p0 *PriceValue, p1 *PriceValue) bool

func CommitCurrenciesFilterFunc(currencies map[string]bool) CommitFilterFunc {
	return func(p0 *PriceValue, p1 *PriceValue) bool {
		if currencies[p1.Currency] {
			return true
		}
		return false
	}
}
func CommitSymbolsFilterFunc(symbols map[string]bool) CommitFilterFunc {
	return func(p0 *PriceValue, p1 *PriceValue) bool {
		if symbols[p1.Symbol] {
			return true
		}
		return false
	}
}
func CommitPricePercentDiffFilterFinc(diffs config.Diffs) CommitFilterFunc {
	return func(p0 *PriceValue, p1 *PriceValue) bool {
		if p0 == nil {
			return true
		}
		diffPercent := diffs.Diff(p1.Symbol)
		if diffPercent >= 100 || diffPercent <= 0 {
			return true
		}
		deltaDiff := p0.Price * float64(diffPercent) / 100
		diff := p1.Price - p0.Price
		if diff < 0 {
			diff *= -1
		}
		if diff >= deltaDiff {
			return true
		}
		return false
	}
}

func CommitTimestampFilterFunc() CommitFilterFunc {
	return func(p0 *PriceValue, p1 *PriceValue) bool {
		if p0 == nil {
			return true
		}
		if p0.TimeStamp.After(p1.TimeStamp) {
			return false
		}
		return true
	}
}
