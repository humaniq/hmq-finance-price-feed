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

type CommitFilterFunc func(p0 *Price, p1 *Price) bool

func CommitCurrenciesFilterFunc(currencies []string) CommitFilterFunc {
	mapper := make(map[string]bool)
	for _, currency := range currencies {
		mapper[currency] = true
	}
	return func(p0 *Price, p1 *Price) bool {
		if mapper[p1.Currency] {
			return true
		}
		return false
	}
}
func CommitSymbolsFilterFunc(symbols []string) CommitFilterFunc {
	mapper := make(map[string]bool)
	for _, symbol := range symbols {
		mapper[symbol] = true
	}
	return func(p0 *Price, p1 *Price) bool {
		if mapper[p1.Symbol] {
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
