package api

import (
	"context"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"net/http"
	"strings"
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
		withHistory := false
		if histQuery := r.URL.Query().Get("history"); histQuery != "" {
			withHistory = true
		}

		resultMap := make(map[string][]PriceRecord)

		for _, currency := range currencies {
			prices, err := backend.GetPrices(ctx, symbols, currencies, withHistory)
			if err != nil {
				if errors.Is(err, app.ErrNotFound) {
					logger.Warn(ctx, "[API] prices not found for %s", currency)
					continue
				}
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
						if withHistory {
							for _, hv := range val.History {
								priceRecord.History = append(priceRecord.History, PriceHistoryRecord{
									TimeStamp: hv.TimeStamp,
									Price:     hv.Value,
								})
							}
						}
						list = append(list, priceRecord)
					}
				}
				resultMap[symbol] = list
			}
		}

		httpext.JSON(w, httpapi.NewOkResponse().WithPayload(resultMap))
	}
}
