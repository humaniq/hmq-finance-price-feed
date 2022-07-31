package svc

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type FeedItem struct {
	records []*PriceRecord
}

func (fi *FeedItem) Filter(ctx context.Context, fn PriceFilterFunc) {
	records := make([]*PriceRecord, 0, len(fi.records))
	for _, record := range fi.records {
		logger.Info(ctx, "record filter: %+v", record)
		if fn(ctx, record) {
			records = append(records, record)
		} else {
			logger.Info(ctx, "filtered out: %+v", record)
		}
	}
	fi.records = records
}

type FeedProvider interface {
	Provide(ctx context.Context, feed chan<- *FeedItem) error
}
type FeedConsumer interface {
	Consume(ctx context.Context, feed <-chan *FeedItem) error
}

type PriceFilterFunc func(ctx context.Context, record *PriceRecord) bool

type PriceFeedHandler struct {
	back    *PriceSvc
	in      chan *FeedItem
	outs    []chan *FeedItem
	filters []PriceFilterFunc
}

func NewPriceFeedHandler(state *PriceSvc) *PriceFeedHandler {
	return &PriceFeedHandler{
		back: state,
		in:   make(chan *FeedItem),
	}
}
func (h *PriceFeedHandler) In() chan<- *FeedItem {
	return h.in
}
func (h *PriceFeedHandler) GetConsumerChan() <-chan *FeedItem {
	c := make(chan *FeedItem)
	h.outs = append(h.outs, c)
	return c
}
func (h *PriceFeedHandler) AddFilterFunc(fn PriceFilterFunc) {
	h.filters = append(h.filters, fn)
}

func (h *PriceFeedHandler) Start() error {
	go func(c chan *FeedItem) {
		ctx := context.Background()
		logger.Info(ctx, "PriceFeedHandler start")
		defer func() {
			for _, out := range h.outs {
				close(out)
			}
		}()
		for item := range c {
			logger.Debug(ctx, "feed_old: %+v", item)
			records := item.records
			for _, filter := range h.filters {
				logger.Debug(ctx, "filtering %+v", records)
				item.Filter(ctx, filter)
			}
			for _, record := range records {
				if err := h.back.SetSymbolPrice(ctx, record.Symbol, record.Currency, record.Price, record.Symbol); err != nil {
					logger.Error(ctx, "error writing to DS: %s", err.Error())
				}
				logger.Info(ctx, "record written: %+v", record)
			}
			for _, out := range h.outs {
				out <- item
			}
		}
	}(h.in)
	return nil
}
func (h *PriceFeedHandler) Stop() error {
	close(h.in)
	return nil
}
