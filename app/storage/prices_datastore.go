package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

type dsPricesRecord struct {
	Key       string         `datastore:"key"`
	TimeStamp time.Time      `datastore:"timeStamp"`
	Prices    []*state.Price `datastore:"prices"`
}

func (r *dsPricesRecord) ToState() *state.Prices {
	prices := state.NewPrices(r.Key)
	for _, price := range r.Prices {
		prices.Commit(price)
	}
	prices.Stage()
	return prices
}
func dsPricesRecordFromState(value *state.Prices) *dsPricesRecord {
	return &dsPricesRecord{
		Key:       value.Key(),
		TimeStamp: time.Now(),
		Prices:    value.Prices(),
	}
}

type PricesDS struct {
	client *gds.Client
}

func NewPricesDS(client *gds.Client) *PricesDS {
	return &PricesDS{client: client}
}

func (ds *PricesDS) SavePrices(ctx context.Context, key string, value *state.Prices) error {
	if err := dsWritePrices(ctx, ds.client, key, dsPricesRecordFromState(value)); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func (ds *PricesDS) LoadPrices(ctx context.Context, key string) (*state.Prices, error) {
	pricesDS, err := dsReadPrices(ctx, ds.client, key)
	if err != nil {
		if errors.Is(err, gds.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return pricesDS.ToState(), nil
}

func dsWritePrices(ctx context.Context, ds *gds.Client, key string, record *dsPricesRecord) error {
	if err := ds.Write(ctx, toPricesDSKey(key), record); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func dsReadPrices(ctx context.Context, ds *gds.Client, key string) (*dsPricesRecord, error) {
	var prices dsPricesRecord
	if err := ds.Read(ctx, toPricesDSKey(key), &prices); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return &prices, nil
}
func toPricesDSKey(symbol string) string {
	return symbol
}
