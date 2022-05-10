package state

import "time"

type Price struct {
	TimeStamp time.Time
	Source    string
	Symbol    string
	Currency  string
	Price     float64
}

func NewPrice(source string, symbol string, currency string, price float64, timeStamp time.Time) *Price {
	return &Price{
		TimeStamp: timeStamp,
		Source:    source,
		Symbol:    symbol,
		Currency:  currency,
		Price:     price,
	}
}

type Prices struct {
	currency      string
	values        map[string]*Price
	changed       []string
	commitFilters []CommitFilterFunc
}

func NewPrices(currency string) *Prices {
	return &Prices{
		currency: currency,
		values:   make(map[string]*Price),
		changed:  []string{},
	}
}
func (ps *Prices) Key() string {
	return ps.currency
}
func (ps *Prices) WithCommitFilters(fn ...CommitFilterFunc) *Prices {
	ps.commitFilters = append(ps.commitFilters, fn...)
	return ps
}
func (ps *Prices) Commit(price *Price) bool {
	current := ps.values[price.Symbol]
	for _, filter := range ps.commitFilters {
		if !filter(current, price) {
			return false
		}
	}
	ps.values[price.Symbol] = price
	ps.changed = append(ps.changed, price.Symbol)
	return true
}
func (ps *Prices) Stage() {
	ps.changed = []string{}
}
func (ps *Prices) Changes() []*Price {
	result := make([]*Price, 0, len(ps.changed))
	for _, change := range ps.changed {
		result = append(result, ps.values[change])
	}
	return result
}
func (ps *Prices) Values() map[string]*Price {
	return ps.values
}
func (ps *Prices) Prices() []*Price {
	result := make([]*Price, 0, len(ps.values))
	for _, val := range ps.values {
		result = append(result, val)
	}
	return result
}
func (ps *Prices) Estimate(symbol string, currency string) *Price {
	symbolPrice, found := ps.values[symbol]
	if !found {
		return nil
	}
	currencyPrice, found := ps.values[currency]
	if !found {
		return nil
	}
	return &Price{
		TimeStamp: time.Now(),
		Source:    "estimate",
		Symbol:    symbol,
		Currency:  currency,
		Price:     symbolPrice.Price / currencyPrice.Price,
	}
}

type CommitFilterFunc func(p0 *Price, p1 *Price) bool

func CommitCurrenciesFilterFunc(currencies map[string]bool) CommitFilterFunc {
	return func(p0 *Price, p1 *Price) bool {
		if currencies[p1.Currency] {
			return true
		}
		return false
	}
}
func CommitSymbolsFilterFunc(symbols map[string]bool) CommitFilterFunc {
	return func(p0 *Price, p1 *Price) bool {
		if symbols[p1.Symbol] {
			return true
		}
		return false
	}
}
func CommitPricePercentDiffFilterFinc(diffs map[string]int) CommitFilterFunc {
	return func(p0 *Price, p1 *Price) bool {
		if p0 == nil {
			return true
		}
		diffPercent, found := diffs[p1.Symbol]
		if !found || diffPercent >= 100 {
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
	return func(p0 *Price, p1 *Price) bool {
		if p0 == nil {
			return true
		}
		if p0.TimeStamp.After(p1.TimeStamp) {
			return false
		}
		return true
	}
}
