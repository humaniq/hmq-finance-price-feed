package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"time"
)

type StorageWriteWorker struct {
	backend storage.Prices
	items   map[string]*price.Asset
}

func NewStorageWriteWorker(backend storage.Prices) *StorageWriteWorker {
	return &StorageWriteWorker{
		backend: backend,
		items:   make(map[string]*price.Asset),
	}
}

func (sww *StorageWriteWorker) Work(ctx context.Context, values []price.Value) error {
	for _, value := range values {
		item, found := sww.items[value.Currency]
		if !found {
			loadedItem, err := sww.backend.LoadPrices(ctx, value.Currency)
			if err != nil {
				return err
			}
			item = loadedItem
		}
		priceHistory := item.History[value.Symbol]
		recalcHistory := false
		if len(priceHistory) > 0 && priceHistory[0].TimeStamp.Before(time.Now().Add(990*time.Hour)) {
			recalcHistory = true
		}
		priceHistory = priceHistory.AddRecords(recalcHistory, price.HistoryRecord{
			TimeStamp: value.TimeStamp,
			Price:     value.Price,
		})
		writtenValue := price.Asset{
			Name:    value.Currency,
			Prices:  map[string]price.Value{value.Symbol: value},
			History: map[string]price.History{value.Symbol: priceHistory},
		}
		if err := sww.backend.SavePrices(ctx, value.Currency, &writtenValue); err != nil {
			app.Logger().Error(ctx, "error writing value %+v: %s", writtenValue, err)
			return err
		}
		item.Prices[value.Symbol] = value
		item.History[value.Symbol] = priceHistory
		sww.items[value.Currency] = item
	}
	return nil
}
