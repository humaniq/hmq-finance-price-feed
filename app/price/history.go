package price

import (
	"sort"
	"time"
)

type History []HistoryRecord

func (h History) AddRecord(record HistoryRecord) History {
	historyStart := time.Now().Add(-960 * time.Hour)
	year, month, day := historyStart.Date()
	historyStart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)

	newHistory := make(*History, 264, 264)
	for _, hValue := range h {
		if hValue.TimeStamp.Before(historyStart) {
			continue
		}
		if hValue.TimeStamp.After(currentHours) {
			newHistory = newHistory.AddRecord(record)
			continue
		}

		if hValue.TimeStamp.After(latestHourly) {

			continue
		}

	}
	return newHistory
}
func historyIndex(start time.Time, current time.Time) (time.Time, int) {
	hoursSinceStart := current.Sub(start) / time.Hour
	if hoursSinceStart <= 912 {
		index := int(hoursSinceStart / 4)
		timeStamp := start.Add(time.Hour * hoursSinceStart)
		return timeStamp, index
	}

}

func (h History) Sort() {
	sort.Slice(h, func(i, j int) bool {
		return h[i].TimeStamp.Before(h[j].TimeStamp)
	})
}

type HistoryRecord struct {
	TimeStamp time.Time `json:"time"`
	Price     float64   `json:"price"`
}
