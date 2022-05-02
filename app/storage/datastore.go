package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

type dsPricesRecord struct {
	Source    string             `datastore:"source"`
	TimeStamp time.Time          `datastore:"timeStamp"`
	Prices    map[string]float64 `datastore:"prices"`
}
type dsPrices struct {
	Symbol  string           `datastore:"symbol"`
	Records []dsPricesRecord `datastore:"records"`
}

type Datastore struct {
	client       *gds.Client
	recordsCount int
}

func NewDatastore(client *gds.Client, historyCount int) *Datastore {
	return &Datastore{client: client, recordsCount: historyCount}
}

func (ds *Datastore) SetSymbolPrices(ctx context.Context, symbol string, source string, timeStamp time.Time, prices map[string]float64) error {
	record, err := dsReadPrices(ctx, ds.client, symbol)
	if err != nil && !errors.Is(err, gds.ErrNotFound) {
		return fmt.Errorf("%w: %s", ErrReading, err)
	}
	records := make([]dsPricesRecord, 0, ds.recordsCount+1)
	records = append(records, dsPricesRecord{
		Source:    source,
		TimeStamp: timeStamp,
		Prices:    prices,
	})
	if record != nil {
		for index, rec := range record.Records {
			if index >= ds.recordsCount {
				break
			}
			records = append(records, rec)
		}
	}
	if err := dsWritePrices(ctx, ds.client, &dsPrices{
		Symbol:  symbol,
		Records: records,
	}); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func (ds *Datastore) GetLatestSymbolPrices(ctx context.Context, symbol string, currencies []string) (*Prices, error) {
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
	for _, pricesRecord := range pricesDS.Records {
		for currency, price := range pricesRecord.Prices {
			prices.PutPrice(
				currency,
				&Price{
					Source:    pricesRecord.Source,
					Price:     price,
					TimeStamp: pricesRecord.TimeStamp,
				},
				filterCurrencies,
			)
		}
	}
	return prices, nil
}

func dsWritePrices(ctx context.Context, ds *gds.Client, record *dsPrices) error {
	if err := ds.Write(ctx, toPricesKey(record.Symbol), record); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func dsReadPrices(ctx context.Context, ds *gds.Client, symbol string) (*dsPrices, error) {
	var prices dsPrices
	if err := ds.Read(ctx, toPricesKey(symbol), &prices); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return nil, nil
}
func toPricesKey(symbol string) string {
	return symbol
}
