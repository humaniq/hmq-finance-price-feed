package price

import (
	"log"
	"sort"
	"time"
)

type History []HistoryRecord

func (h History) AddRecords(recalculate bool, records ...HistoryRecord) History {
	historyRecords := append(h, records...)
	historyRecords.Sort()
	if !recalculate {
		return historyRecords
	}
	historyStart := time.Now().Add(-960 * time.Hour)
	year, month, day := historyStart.Date()
	historyStart = time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
	latestHours := time.Now().Add(-24 * time.Hour)
	latestWeek := latestHours.Add(-7 * 24 * time.Hour)
	monthHistory := make(History, 132, 132)
	weekHistory := make(History, 56, 56)
	var latestHistory []HistoryRecord
	for _, hValue := range historyRecords {
		if hValue.TimeStamp.Before(historyStart) {
			continue
		}
		if hValue.TimeStamp.After(latestHours) {
			latestHistory = append(latestHistory, hValue)
			continue
		}
		if hValue.TimeStamp.After(latestWeek) {
			ts, index := historyWeekIndex(latestWeek, hValue.TimeStamp)
			weekHistory[index] = HistoryRecord{
				TimeStamp: ts,
				Price:     hValue.Price,
			}
			continue
		}
		ts, index := historyMonthIndex(historyStart, hValue.TimeStamp)
		monthHistory[index] = HistoryRecord{
			TimeStamp: ts,
			Price:     hValue.Price,
		}
	}
	historyRecords = make(History, 0, len(latestHistory)+len(monthHistory)+len(weekHistory))
	for _, hr := range monthHistory {
		if hr.TimeStamp.IsZero() {
			continue
		}
		historyRecords = append(historyRecords, hr)
	}
	for _, hr := range weekHistory {
		if hr.TimeStamp.IsZero() {
			continue
		}
		historyRecords = append(historyRecords, hr)
	}
	for _, hr := range latestHistory {
		if hr.TimeStamp.IsZero() {
			continue
		}
		historyRecords = append(historyRecords, hr)
	}
	return historyRecords
}
func historyWeekIndex(start time.Time, current time.Time) (time.Time, int) {
	index := current.Sub(start) / time.Hour / 3
	log.Printf("weekIndex: %v=>%v = %d\n", start, current, index)
	return start.Add(time.Hour * 3 * index), int(index)
}
func historyMonthIndex(start time.Time, current time.Time) (time.Time, int) {
	index := current.Sub(start) / time.Hour / 6
	log.Printf("monthIndex: %v=>%v = %d\n", start, current, index)
	return start.Add(time.Hour * 6 * index), int(index)
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
