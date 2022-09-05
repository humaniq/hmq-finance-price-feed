package prices

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
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
func (cs *ConsumerState) Work(ctx context.Context, values []price.Value) error {
	for _, value := range values {
		cs.stateMap[cs.keyFn(value)] = value
	}
	return nil
}
func (cs *ConsumerState) EnrichFunc(symbols []string, currency string, assets []string) func(ctx context.Context, value price.Value) []price.Value {
	symbolToEnrich := make(map[string]bool)
	for _, symbol := range symbols {
		symbolToEnrich[symbol] = true
	}
	return func(ctx context.Context, value price.Value) []price.Value {
		result := make([]price.Value, 0, len(assets))
		result = append(result, value)
		if value.Currency != currency {
			return result
		}
		for _, asset := range assets {
			if asset == value.Currency {
				continue
			}
			mappingValue := price.Value{
				Symbol:   asset,
				Currency: currency,
			}
			mapValue, found := cs.stateMap[cs.keyFn(mappingValue)]
			if !found {
				continue
			}
			result = append(result, price.Value{
				TimeStamp: value.TimeStamp,
				Source:    value.Source,
				Symbol:    value.Symbol,
				Currency:  asset,
				Price:     value.Price / mapValue.Price,
			})
		}
		return result
	}
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

type ThresholdTimeDeltas struct {
	mapper       map[string]time.Duration
	defaultDelta time.Duration
}

func NewThresholdTimeDeltas(deltas config.Thresholds) *ThresholdTimeDeltas {
	mapper := make(map[string]time.Duration)
	for _, threshold := range deltas.Custom {
		mapper[fmt.Sprintf("%s:%s", threshold.Symbol, threshold.Currency)] = threshold.TimeDelta()
	}
	//app.Logger().Info(context.Background(), "THRESHOLD_TIME_DELTAS: %+v", mapper)
	return &ThresholdTimeDeltas{mapper: mapper, defaultDelta: deltas.Default.TimeDelta()}
}
func (td *ThresholdTimeDeltas) TimeDelta(value price.Value) time.Duration {
	if delta, found := td.mapper[fmt.Sprintf("%s:%s", value.Symbol, value.Currency)]; found {
		return delta
	}
	if delta, found := td.mapper[fmt.Sprintf("%s:", value.Symbol)]; found {
		return delta
	}
	if delta, found := td.mapper[fmt.Sprintf(":%s", value.Currency)]; found {
		return delta
	}
	return td.defaultDelta
}
func (cs *ConsumerState) TimeDeltaThresholdsFunc(deltas config.Thresholds) func(ctx context.Context, value price.Value) bool {
	thresholds := NewThresholdTimeDeltas(deltas)
	return func(ctx context.Context, value price.Value) bool {
		delta := thresholds.TimeDelta(value)
		//app.Logger().Info(context.Background(), "TIME_DELTA: %+v => %v", value, delta)
		if delta == 0 {
			return false
		}
		if currentValue, found := cs.stateMap[cs.keyFn(value)]; found && currentValue.TimeStamp.Add(delta).Before(time.Now()) {
			return true
		}
		return false
	}
}

type ThresholdPercentDeltas struct {
	mapper       map[string]float64
	defaultDelta float64
}

func NewThresholdPercentDeltas(deltas config.Thresholds) *ThresholdPercentDeltas {
	mapper := make(map[string]float64)
	for _, threshold := range deltas.Custom {
		mapper[fmt.Sprintf("%s:%s", threshold.Symbol, threshold.Currency)] = threshold.PercentThreshold
	}
	//app.Logger().Info(context.Background(), "THRESHOLD_PERCENTS: %+v", mapper)
	return &ThresholdPercentDeltas{
		mapper:       mapper,
		defaultDelta: deltas.Default.PercentThreshold,
	}
}
func (pt *ThresholdPercentDeltas) PercentDelta(value price.Value) float64 {
	if delta, found := pt.mapper[fmt.Sprintf("%s:%s", value.Symbol, value.Currency)]; found {
		return delta
	}
	if delta, found := pt.mapper[fmt.Sprintf("%s:", value.Symbol)]; found {
		return delta
	}
	if delta, found := pt.mapper[fmt.Sprintf(":%s", value.Currency)]; found {
		return delta
	}
	return pt.defaultDelta
}
func (cs *ConsumerState) PercentThresholdsFunc(deltas config.Thresholds) func(ctx context.Context, value price.Value) bool {
	thresholds := NewThresholdPercentDeltas(deltas)
	return func(ctx context.Context, value price.Value) bool {
		percent := thresholds.PercentDelta(value)
		//app.Logger().Info(context.Background(), "PERCENT_DELTA: %+v => %v", value, percent)
		currentValue := cs.stateMap[cs.keyFn(value)]
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
