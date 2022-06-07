package state

import (
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type CommitValueFilterFunc func(p0 *price.Value, p1 *price.Value) bool

func CommitValueCurrenciesFilterFunc(currencies map[string]bool) CommitValueFilterFunc {
	return func(p0 *price.Value, p1 *price.Value) bool {
		if currencies[p1.Currency] {
			return true
		}
		return false
	}
}
func CommitValueSymbolsFilterFunc(symbols map[string]bool) CommitValueFilterFunc {
	return func(p0 *price.Value, p1 *price.Value) bool {
		if symbols[p1.Symbol] {
			return true
		}
		return false
	}
}
func CommitValuePriceDiffOrTimestampDiffFilterFunc(diffs config.Diffs, ts config.TSDiffs) CommitValueFilterFunc {
	valueDiffFunc := CommitValuePricePercentDiffFilterFinc(diffs)
	tsForceFunc := CommitValuePriceForceTimeStampFilterFunc(ts)
	return func(p0 *price.Value, p1 *price.Value) bool {
		if valueDiffFunc(p0, p1) {
			return true
		}
		if tsForceFunc(p0, p1) {
			return true
		}
		return false
	}
}
func CommitValuePriceForceTimeStampFilterFunc(ts config.TSDiffs) CommitValueFilterFunc {
	return func(p0 *price.Value, p1 *price.Value) bool {
		if p0 == nil {
			return true
		}
		if p1.TimeStamp.Sub(p0.TimeStamp) > ts.Diff(p0.Symbol) {
			return true
		}
		return false
	}
}
func CommitValuePricePercentDiffFilterFinc(diffs config.Diffs) CommitValueFilterFunc {
	return func(p0 *price.Value, p1 *price.Value) bool {
		if p0 == nil {
			return true
		}
		diffPercent := diffs.Diff(p1.Symbol)
		if diffPercent >= 1000000 || diffPercent <= 0 {
			return true
		}
		deltaDiff := p0.Price * float64(diffPercent) / 1000000
		diff := p1.Price - p0.Price
		if diff < 0 {
			diff *= -1
		}
		if diff >= deltaDiff {
			return true
		}
		return false
	}
}
func CommitValueTimestampFilterFunc() CommitValueFilterFunc {
	return func(p0 *price.Value, p1 *price.Value) bool {
		if p0 == nil {
			return true
		}
		if p0.TimeStamp.After(p1.TimeStamp) {
			return false
		}
		return true
	}
}
