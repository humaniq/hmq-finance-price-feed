package legacy

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
