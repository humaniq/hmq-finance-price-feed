package storage

import (
	"context"
	"time"
)

type PricesRecord struct {
	Symbol    string
	Source    string
	TimeStamp time.Time
	Prices    map[string]float64
}

func NewPricesRecord(symbol string, source string, timeStamp time.Time) *PricesRecord {
	return &PricesRecord{
		Symbol:    symbol,
		Source:    source,
		TimeStamp: timeStamp,
		Prices:    make(map[string]float64),
	}
}

type PriceCommitter interface {
	CommitSymbolPrices(ctx context.Context, symbol string, source string, timeStamp time.Time, prices map[string]float64) error
}
type PriceGetter interface {
	GetSymbolPrices(ctx context.Context, symbol string) (*PricesRecord, error)
}
type Pricer interface {
	PriceCommitter
	PriceGetter
}

//
//type Prices struct {
//	Symbol string            `json:"symbol"`
//	Values map[string]*Price `json:"values"`
//}
//
//func NewPrices(symbol string) *Prices {
//	return &Prices{
//		Symbol: symbol,
//		Values: make(map[string]*Price),
//	}
//}
//func (ps *Prices) WithCurrencies(currencies []string) *Prices {
//	for _, currency := range currencies {
//		ps.Values[currency] = nil
//	}
//	return ps
//}
//func (ps *Prices) PutPrice(currency string, price *Price, filterCurrency bool) {
//	current, exists := ps.Values[currency]
//	if !exists && filterCurrency {
//		return
//	}
//	if current == nil {
//		ps.Values[currency] = price
//		return
//	}
//	if current.TimeStamp.After(price.TimeStamp) {
//		if price.TimeStamp.After(current.PreviousTimeStamp) {
//			current.PreviousPrice = price.Price
//			current.PreviousTimeStamp = price.TimeStamp
//			ps.Values[currency] = current
//			return
//		}
//	}
//	if current.TimeStamp.After(price.PreviousTimeStamp) {
//		price.PreviousPrice = current.PreviousPrice
//		price.PreviousTimeStamp = current.TimeStamp
//	}
//	ps.Values[currency] = price
//}
//
//type Price struct {
//	Source            string    `json:"source"`
//	Price             float64   `json:"price"`
//	TimeStamp         time.Time `json:"timeStamp,omitempty"`
//	PreviousPrice     float64   `json:"previousPrice,omitempty"`
//	PreviousTimeStamp time.Time `json:"-"`
//}
//
//func NewPrice(source string, price float64, timeStamp time.Time) *Price {
//	return &Price{
//		Source:    source,
//		Price:     price,
//		TimeStamp: timeStamp,
//	}
//}
//func (p *Price) WithPreviousPrice(value float64, timeStamp time.Time) *Price {
//	p.PreviousPrice = value
//	p.PreviousTimeStamp = timeStamp
//	return p
//}
