package api

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
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

func GetPricesFunc(state svc.PricesGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		symbols := MustGetStringListFromCtx(ctx, CtxSymbolKey)
		currencies := MustGetStringListFromCtx(ctx, CtxCurrencyKey)
		if len(symbols) == 0 || len(currencies) == 0 {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("invalid symbol/currency mapping"), http.StatusBadRequest)
			return
		}

		prices, err := state.GetPrices(ctx, symbols, currencies)
		if err != nil {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("internal error"), http.StatusInternalServerError)
			return
		}

		resultMap := make(map[string][]PriceRecord)

		for symbol, symbolPrices := range prices {
			records := make([]PriceRecord, 0, len(symbolPrices))
			for currency, price := range symbolPrices {
				records = append(records, PriceRecord{
					Source:    price.Source,
					Currency:  currency,
					Price:     price.Value,
					TimeStamp: price.TimeStamp,
				})
			}
			resultMap[symbol] = records
		}

		httpext.JSON(w, httpapi.NewOkResponse().WithPayload(resultMap))
	}
}
