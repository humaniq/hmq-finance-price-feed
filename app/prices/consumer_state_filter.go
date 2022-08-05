package prices

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"time"
)

type StateKeyFunc func(value price.Value) string

func SymbolCurrencyStateKey(value price.Value) string {
	return fmt.Sprintf("%s-%s", value.Symbol, value.Currency)
}
func SymbolStateKey(value price.Value) string {
	return value.Symbol
}

type ConsumerState struct {
	stateMap map[string]price.Value
	keyFn    StateKeyFunc
}

func NewConsumerState(keyFn StateKeyFunc) *ConsumerState {
	return &ConsumerState{
		stateMap: make(map[string]price.Value),
		keyFn:    keyFn,
	}
}
func (cs *ConsumerState) WithValues(values ...price.Value) *ConsumerState {
	for _, value := range values {
		cs.stateMap[cs.keyFn(value)] = value
	}
	return cs
}

func (cs *ConsumerState) ValueExists(ctx context.Context, value price.Value) bool {
	if _, found := cs.stateMap[cs.keyFn(value)]; found {
		return true
	}
	return false
}
func (cs *ConsumerState) TimeDeltaFunc(delta time.Duration) func(context.Context, price.Value) bool {
	return func(ctx context.Context, value price.Value) bool {
		currentValue := cs.stateMap[cs.keyFn(value)]
		if currentValue.TimeStamp.Add(delta).Before(time.Now()) {
			return true
		}
		return false
	}
}
func (cs *ConsumerState) PercentThresholdFunc(thresholds map[string]float64, defaultThreshold float64, keyFn StateKeyFunc) func(ctx context.Context, value price.Value) bool {
	return func(ctx context.Context, value price.Value) bool {
		currentValue := cs.stateMap[cs.keyFn(value)]
		percent := thresholds[keyFn(value)]
		if percent == 0 {
			percent = defaultThreshold
		}
		thresholdDiff := currentValue.Price * percent / 100
		currentDiff := value.Price - currentValue.Price
		if currentDiff < 0 {
			currentDiff = currentDiff * (-1)
		}
		if currentDiff > thresholdDiff {
			return true
		}
		return false
	}
}
