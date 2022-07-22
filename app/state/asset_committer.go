package state

import (
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type AssetCommitter struct {
	asset   *price.Asset
	changed []price.Value
	filters []CommitValueFilterFunc
}

func NewAssetCommitter(asset *price.Asset) *AssetCommitter {
	return &AssetCommitter{asset: asset}
}
func (ac *AssetCommitter) WithFilters(filters ...CommitValueFilterFunc) *AssetCommitter {
	for _, filter := range filters {
		ac.filters = append(ac.filters, filter)
	}
	return ac
}
func (ac *AssetCommitter) Commit(value price.Value) bool {
	currentValue, found := ac.asset.Prices[value.Symbol]
	if found {
		for _, filter := range ac.filters {
			if !filter(&currentValue, &value) {
				return false
			}
		}
	}
	ac.asset.Prices[value.Symbol] = value
	assetHistory := ac.asset.History[value.Symbol]
	assetHistory = append(assetHistory, price.HistoryRecord{
		TimeStamp: value.TimeStamp,
		Price:     value.Price,
	})
	assetHistory.Sort()
	ac.asset.History[value.Symbol] = assetHistory
	ac.changed = append(ac.changed, value)
	return true
}
func (ac *AssetCommitter) Stage() *price.Asset {
	ac.changed = []price.Value{}
	return ac.asset
}
func (ac *AssetCommitter) Changes() []price.Value {
	return ac.changed
}
