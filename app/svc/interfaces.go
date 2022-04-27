package svc

import (
	"context"
	"errors"
	"time"
)

var ErrNoValue = errors.New("no value")

type PriceRecord struct {
	Source        string    `json:"source"`
	Symbol        string    `json:"symbol"`
	Currency      string    `json:"currency"`
	Price         float64   `json:"price"`
	PreviousPrice float64   `json:"previousPrice,omitempty"`
	TimeStamp     time.Time `json:"timeStamp"`
}

func NewPriceRecord(symbol string, currency string, price float64, source string) *PriceRecord {
	return &PriceRecord{
		Source:    source,
		Symbol:    symbol,
		Currency:  currency,
		Price:     price,
		TimeStamp: time.Now(),
	}
}
func (pr *PriceRecord) WithPreviousPrice(price float64) *PriceRecord {
	pr.PreviousPrice = price
	return pr
}

type PriceGetter interface {
	GetLatestSymbolPrice(ctx context.Context, symbol string, currency string) (*PriceRecord, error)
}
type PriceSetter interface {
	SetSymbolPrice(ctx context.Context, price *PriceRecord) error
}

type Pricer interface {
	PriceSetter
	PriceGetter
}

type PriceFilterFunc func(ctx context.Context, record *PriceRecord) bool

type FeedItem struct {
	records []*PriceRecord
}

func (fi *FeedItem) Filter(ctx context.Context, fn PriceFilterFunc) {
	records := make([]*PriceRecord, 0, len(fi.records))
	for _, record := range fi.records {
		if fn(ctx, record) {
			records = append(records, record)
		}
	}
	fi.records = records
}

type FeedReleaser interface {
	Lease() chan<- *FeedItem
	Release()
}
type AsyncRunner interface {
	Start() error
	Stop() error
	WaitForDone()
}
type AsyncConsumer interface {
	In() chan<- *FeedItem
	AsyncRunner
	FeedReleaser
}

type FeedProvider interface {
	Provide(ctx context.Context, feed chan<- *FeedItem) error
}
type FeedConsumer interface {
	Consume(ctx context.Context, feed <-chan *FeedItem) error
}
