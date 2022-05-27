package api

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"net/http"
	"sort"
	"strconv"
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

		precisionRequest := r.URL.Query().Get("historyPrecision")
		precisionCount, err := strconv.Atoi(precisionRequest)
		if err != nil {
			precisionCount = 48
		}

		resultMap := make(map[string]map[string]PriceRecord)

		prices, err := backend.GetPrices(ctx, symbols, currencies, withHistory)
		if err != nil {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("error getting prices"), http.StatusInternalServerError)
			return
		}
		for _, symbol := range symbols {
			list, found := resultMap[symbol]
			if !found {
				list = make(map[string]PriceRecord)
			}
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
						case "custom":
							sinceRequest := r.URL.Query().Get("sinceTimestamp")
							since := time.Now().Add(time.Hour * 24 * (-7))
							sinceNumber, err := strconv.ParseInt(sinceRequest, 10, 64)
							if err == nil {
								since = time.Unix(sinceNumber, 0)
							}
							priceRecord.History = buildHistoryChart(since, precisionCount, val.History)
							break
						case "year":
							priceRecord.History = buildHistoryChart(time.Now().Add(-8760*time.Hour), precisionCount, val.History)
							break
						case "month":
							priceRecord.History = buildHistoryChart(time.Now().Add(-720*time.Hour), precisionCount, val.History)
							break
						case "week":
							priceRecord.History = buildHistoryChart(time.Now().Add(-168*time.Hour), precisionCount, val.History)
							break
						case "day":
							priceRecord.History = buildHistoryChart(time.Now().Add(-24*time.Hour), precisionCount, val.History)
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
	period := now.Sub(since) / time.Duration(granularity-1)
	for index := 0; index < granularity; index++ {
		result = append(result, &PriceHistoryRecord{
			TimeStamp: since.Add(period * time.Duration(index)),
			Price:     0,
		})
	}
	estimator := make([]*PriceHistoryRecord, 0, granularity+len(records))
	estimator = append(estimator, result...)
	for _, record := range records {
		estimator = append(estimator, &PriceHistoryRecord{
			TimeStamp: record.TimeStamp,
			Price:     record.Value,
		})
	}
	sort.Slice(estimator, func(i, j int) bool {
		return estimator[i].TimeStamp.Before(estimator[j].TimeStamp)
	})
	currentPrice := float64(0)
	for _, e := range estimator {
		logger.Info(context.Background(), "ESTIMATING: %+v", e)
		if e.Price == 0 {
			e.Price = currentPrice
		} else {
			currentPrice = e.Price
		}
	}
	return result
}
