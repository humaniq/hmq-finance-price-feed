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

type dsPriceAsset struct {
	Currency  string            `datastore:"currency"`
	Symbol    string            `datastore:"symbol"`
	TimeStamp time.Time         `datastore:"timestamp,noindex"`
	Price     price.Value       `datastore:"price,noindex"`
	History   []dsHistoryRecord `datastore:"history,noindex"`
}

type PricesDSv2 struct {
	client *gds.Client
}

func NewPricesDSv2(client *gds.Client) *PricesDSv2 {
	return &PricesDSv2{client: client}
}

func (pdsv2 *PricesDSv2) SavePrices(ctx context.Context, key string, value *price.Asset) error {
	assets := make([]*dsPriceAsset, 0, len(value.Prices))
	for symbol, val := range value.Prices {
		asset := dsPriceAsset{
			Currency:  val.Currency,
			Symbol:    val.Symbol,
			TimeStamp: val.TimeStamp,
			Price:     val,
		}
		history, found := value.History[symbol]
		if found && len(history) > 0 {
			asset.History = make([]dsHistoryRecord, 0, len(history))
			for _, record := range history {
				asset.History = append(asset.History, dsHistoryRecord{
					Symbol:    symbol,
					TimeStamp: record.TimeStamp,
					Price:     record.Price,
				})
			}
		}
		assets = append(assets, &asset)
	}
	return dsV2SavePrices(ctx, pdsv2.client, assets)
}
func (pdsv2 *PricesDSv2) LoadPrices(ctx context.Context, key string) (*price.Asset, error) {
	dsAssets, err := dsV2ReadPrices(ctx, pdsv2.client, key)
	if err != nil {
		if errors.Is(err, gds.ErrNotFound) {
			return nil, app.ErrNotFound
		}
		return nil, err
	}
	asset := price.NewAsset(key)
	for _, dsAsset := range dsAssets {
		priceVal := dsAsset.Price
		asset.Prices[dsAsset.Symbol] = price.Value{
			TimeStamp: priceVal.TimeStamp,
			Source:    priceVal.Source,
			Symbol:    priceVal.Symbol,
			Currency:  priceVal.Currency,
			Price:     priceVal.Price,
		}
		history := make([]price.HistoryRecord, 0, len(dsAsset.History))
		for _, dsHistoryVal := range dsAsset.History {
			history = append(history, price.HistoryRecord{
				TimeStamp: dsHistoryVal.TimeStamp,
				Price:     dsHistoryVal.Price,
			})
		}
		asset.History[dsAsset.Symbol] = history
	}
	return asset, nil
}

func dsV2SavePrices(ctx context.Context, ds *gds.Client, records []*dsPriceAsset) error {
	valuesMap := make(map[string]interface{})
	for _, value := range records {
		valuesMap[toPricesDSv2Key(value.Currency, value.Symbol)] = value
	}
	return gds.WriteMultiple(ctx, ds, valuesMap)
}
func dsV2ReadPrices(ctx context.Context, ds *gds.Client, currency string) ([]dsPriceAsset, error) {
	return gds.ReadMultipleByFilters[dsPriceAsset](ctx, ds, []gds.Filter{{
		Str:   "currency =",
		Value: currency,
	}})
}
func toPricesDSv2Key(currency string, symbol string) string {
	return fmt.Sprintf("%s_%s", currency, symbol)
}
