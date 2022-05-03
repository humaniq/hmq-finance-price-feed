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
func (ds *DatastorePricer) GetSymbolPrices(ctx context.Context, symbol string, currencies []string) (*Prices, error) {
	pricesDS, err := dsReadPrices(ctx, ds.client, symbol)
	if err != nil {
		if errors.Is(err, gds.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	prices := NewPrices(pricesDS.Symbol)
	filterCurrencies := false
	if currencies != nil {
		prices = prices.WithCurrencies(currencies)
		filterCurrencies = true
	}
	for currency, price := range pricesDS.Prices {
		prices.PutPrice(currency, NewPrice(pricesDS.Source, price, pricesDS.TimeStamp), filterCurrencies)
	}
	return prices, nil
}

func dsWritePrices(ctx context.Context, ds *gds.Client, record *dsPricesRecord) error {
	if err := ds.Write(ctx, toPricesKey(record.Symbol), record); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func dsReadPrices(ctx context.Context, ds *gds.Client, symbol string) (*dsPricesRecord, error) {
	var prices dsPricesRecord
	if err := ds.Read(ctx, toPricesKey(symbol), &prices); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return &prices, nil
}
func toPricesKey(symbol string) string {
	return symbol
}
