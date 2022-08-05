package price

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
	ar.Price = value
	ar.History = ar.History.AddRecord(HistoryRecord{
		TimeStamp: value.TimeStamp,
		Price:     value.Price,
	})
}
