package api

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
)

const CtxSymbolKey = "symbol"
const CtxCurrencyKey = "currency"

func MustHaveStringListInQueryOrDefaultsMiddlewareFunc(queryKey string, ctxKey string, caseCast int, delimiter string, defaults []string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			stringValue := r.URL.Query().Get(queryKey)
			if stringValue == "" {
				if defaults != nil && len(defaults) > 0 {
					next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, ctxKey, defaults)))
					return
				}
				httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("query is empty"), http.StatusBadRequest)
				return
			}
			list := strings.Split(stringValue, delimiter)
			if caseCast != httpapi.CaseSensitive {
				castList := make([]string, 0, len(list))
				for _, item := range list {
					if caseCast == httpapi.CaseToUpper {
						castList = append(castList, strings.ToUpper(item))
					}
					if caseCast == httpapi.CaseToLower {
						castList = append(castList, strings.ToLower(item))
					}
				}
				list = castList
			}
			next.ServeHTTP(w, r.WithContext(context.WithValue(ctx, ctxKey, list)))
		})
	}
}
func MustGetStringListFromCtx(ctx context.Context, key string) []string {
	return ctx.Value(key).([]string)
}

func GetPricesFunc(backend svc.PricesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		symbols := MustGetStringListFromCtx(ctx, CtxSymbolKey)
		currencies := MustGetStringListFromCtx(ctx, CtxCurrencyKey)
		if len(symbols) == 0 || len(currencies) == 0 {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("invalid symbol/currency mapping"), http.StatusBadRequest)
			return
		}
		history := r.URL.Query().Get("history")
		withHistory := false
		if history != "" {
			withHistory = true
		}

		resultMap := make(map[string]map[string]PriceRecord)

		prices, err := backend.GetPrices(ctx, symbols, currencies, withHistory)
		if err != nil {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("error getting prices"), http.StatusInternalServerError)
			return
		}
		for _, symbol := range symbols {
			list := resultMap[symbol]
			value, found := prices[symbol]
			if found {
				for key, val := range value {
					priceRecord := PriceRecord{
						Source:    val.Source,
						Currency:  key,
						TimeStamp: val.TimeStamp,
						Price:     val.Value,
					}
					if history != "" {
						switch history {
						case "year":
							priceRecord.History = buildHistoryChart(time.Now().Add(-8760*time.Hour), 48, val.History)
							break
						case "month":
							priceRecord.History = buildHistoryChart(time.Now().Add(-720*time.Hour), 48, val.History)
							break
						case "week":
							priceRecord.History = buildHistoryChart(time.Now().Add(-168*time.Hour), 48, val.History)
							break
						case "day":
							priceRecord.History = buildHistoryChart(time.Now().Add(-24*time.Hour), 48, val.History)
							break
						default:
							for _, hv := range val.History {
								priceRecord.History = append(priceRecord.History, &PriceHistoryRecord{
									TimeStamp: hv.TimeStamp,
									Price:     hv.Value,
								})
							}
						}
					}
					list[priceRecord.Currency] = priceRecord
				}
			}
			resultMap[symbol] = list
		}

		httpext.JSON(w, httpapi.NewOkResponse().WithPayload(resultMap))
	}
}

func buildHistoryChart(since time.Time, granularity int, records []svc.SymbolPricesHistory) []*PriceHistoryRecord {
	now := time.Now()
	result := make([]*PriceHistoryRecord, 0, granularity)
	period := now.Sub(since) / time.Duration(granularity)
	for index := 0; index < granularity; index++ {
		result = append(result, &PriceHistoryRecord{
			TimeStamp: since.Add(period * time.Duration(index)),
			Price:     0,
		})
	}
	cursor := 0
	for _, record := range records {
		if cursor >= granularity {
			break
		}
		cursorValue := result[cursor]
		if record.TimeStamp.Before(cursorValue.TimeStamp) {
			cursorValue.Price = record.Value
			continue
		}
		for cursor < granularity {
			cursor++
			cursorValue = result[cursor]
			cursorValue.Price = record.Value
			if record.TimeStamp.Before(cursorValue.TimeStamp) {
				break
			}
		}
	}
	if cursor < granularity-1 {
		latestValue := result[cursor].Price
		for cursor < granularity {
			result[cursor].Price = latestValue
			cursor++
		}
	}
	return result
}
