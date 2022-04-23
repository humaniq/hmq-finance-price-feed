package svc

import (
	"context"
)

type FeedItem struct {
	records []*PriceRecord
}

type FeedHandler interface {
	HandleFeed(ctx context.Context, feed chan *FeedItem)
}
type FeedHandlerFunc func(ctx context.Context, feed chan *FeedItem)

func (fhf FeedHandlerFunc) HandleFeed(ctx context.Context, feed chan *FeedItem) {
	fhf(ctx, feed)
}

type FeedProvider interface {
	Provide(ctx context.Context, feed chan<- *FeedItem) error
}
type FeedConsumer interface {
	Consume(ctx context.Context, feed <-chan *FeedItem) error
}
