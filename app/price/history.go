package price

import (
	"sort"
	"time"
)

type History []HistoryRecord

func (h History) Sort() {
	sort.Slice(h, func(i, j int) bool {
		return h[i].TimeStamp.Before(h[j].TimeStamp)
	})
}

type HistoryRecord struct {
	TimeStamp time.Time `json:"time"`
	Price     float64   `json:"price"`
}
