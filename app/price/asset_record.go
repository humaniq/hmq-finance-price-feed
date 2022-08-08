package price

import "time"

type AssetRecord struct {
	Price   Value
	History History
}

func NewAssetRecord(value Value) *AssetRecord {
	return &AssetRecord{
		Price: value,
	}
}
func (ar *AssetRecord) WithHistory(history History) *AssetRecord {
	ar.History = history
	return ar
}
func (ar *AssetRecord) Commit(value Value) {
	recalculate := false
	if len(ar.History) > 0 && value.TimeStamp.Add(time.Hour*990).After(ar.History[0].TimeStamp) {
		recalculate = true
	}
	ar.Price = value
	ar.History = ar.History.AddRecords(recalculate, HistoryRecord{
		TimeStamp: value.TimeStamp,
		Price:     value.Price,
	})
}
