package storage

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
)

type dsPricesAsset struct {
	Key       string            `datastore:"key"`
	TimeStamp time.Time         `datastore:"timeStamp,noindex"`
	Prices    []price.Value     `datastore:"prices,noindex"`
	History   []dsHistoryRecord `datastore:"history,noindex"`
}
type dsHistoryRecord struct {
	Symbol    string    `datastore:"key,noindex"`
	TimeStamp time.Time `datastore:"timestamp,noindex"`
	Price     float64   `datastore:"price,noindex"`
}

func (r *dsPricesAsset) ToAsset() *price.Asset {
	prices := price.NewAsset(r.Key)
	names := make([]string, 0, len(r.Prices))
	for _, priceValue := range r.Prices {
		names = append(names, priceValue.Symbol)
		prices.Prices[priceValue.Symbol] = priceValue
	}
	for _, historyRecord := range r.History {
		item, found := prices.History[historyRecord.Symbol]
		if !found {
			item = price.History{}
		}
		item = append(item, price.HistoryRecord{
			TimeStamp: historyRecord.TimeStamp,
			Price:     historyRecord.Price,
		})
		prices.History[historyRecord.Symbol] = item
	}
	return prices
}
func dsPricesRecordFromAsset(value *price.Asset) *dsPricesAsset {
	values := make([]price.Value, 0, len(value.Prices))
	for _, priceValue := range value.Prices {
		values = append(values, priceValue)
	}
	var history []dsHistoryRecord
	for symbol, priceHistory := range value.History {
		for _, historyItem := range priceHistory {
			history = append(history, dsHistoryRecord{
				Symbol:    symbol,
				TimeStamp: historyItem.TimeStamp,
				Price:     historyItem.Price,
			})
		}
	}
	return &dsPricesAsset{
		Key:       value.Name,
		TimeStamp: time.Now(),
		Prices:    values,
		History:   history,
	}
}

type PricesDS struct {
	client *gds.Client
}

func NewPricesDS(client *gds.Client) *PricesDS {
	return &PricesDS{client: client}
}

func (ds *PricesDS) SavePrices(ctx context.Context, key string, value *price.Asset) error {
	if err := dsWritePrices(ctx, ds.client, key, dsPricesRecordFromAsset(value)); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func (ds *PricesDS) LoadPrices(ctx context.Context, key string) (*price.Asset, error) {
	pricesDS, err := dsReadPrices(ctx, ds.client, key)
	if err != nil {
		return nil, err
	}
	return pricesDS.ToAsset(), nil
}

func dsWritePrices(ctx context.Context, ds *gds.Client, key string, record *dsPricesAsset) error {
	if err := ds.Write(ctx, toPricesDSKey(key), record); err != nil {
		return fmt.Errorf("%w: %s", ErrWriting, err)
	}
	return nil
}
func dsReadPrices(ctx context.Context, ds *gds.Client, key string) (*dsPricesAsset, error) {
	var prices dsPricesAsset
	if err := ds.Read(ctx, toPricesDSKey(key), &prices); err != nil {
		if errors.Is(err, gds.ErrNotFound) {
			return nil, app.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %s", ErrReading, err)
	}
	return &prices, nil
}
func toPricesDSKey(symbol string) string {
	return symbol
}
