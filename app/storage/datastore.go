package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

type dsPricesRecord struct {
	Symbol    string             `datastore:"symbol"`
	Source    string             `datastore:"source"`
	TimeStamp time.Time          `datastore:"timeStamp"`
	Prices    map[string]float64 `datastore:"prices"`
}

type DatastorePricer struct {
	client *gds.Client
}

func NewDatastorePricer(client *gds.Client) *DatastorePricer {
	return &DatastorePricer{client: client}
}
func (ds *DatastorePricer) CommitSymbolPrices(ctx context.Context, symbol string, source string, timeStamp time.Time, prices map[string]float64) error {
	record, err := dsReadPrices(ctx, ds.client, symbol)
	if err != nil && !errors.Is(err, gds.ErrNotFound) {
		return fmt.Errorf("%w: %s", ErrReading, err)
	}
	if record != nil && record.TimeStamp.After(timeStamp) {
		return ErrTooLate
	}
	if err := dsWritePrices(ctx, ds.client, &dsPricesRecord{
		Symbol:    symbol,
		Source:    source,
		TimeStamp: timeStamp,
		Prices:    prices,
	}); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func (ds *DatastorePricer) GetSymbolPrices(ctx context.Context, symbol string) (*PricesRecord, error) {
	pricesDS, err := dsReadPrices(ctx, ds.client, symbol)
	if err != nil {
		if errors.Is(err, gds.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	prices := NewPricesRecord(symbol, pricesDS.Source, pricesDS.TimeStamp)
	for currency, price := range pricesDS.Prices {
		prices.Prices[currency] = price
	}
	return prices, nil
}

func dsWritePrices(ctx context.Context, ds *gds.Client, record *dsPricesRecord) error {
	if err := ds.Write(ctx, toPricesDSKey(record.Symbol), record); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func dsReadPrices(ctx context.Context, ds *gds.Client, symbol string) (*dsPricesRecord, error) {
	var prices dsPricesRecord
	if err := ds.Read(ctx, toPricesDSKey(symbol), &prices); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return &prices, nil
}
func toPricesDSKey(symbol string) string {
	return symbol
}
