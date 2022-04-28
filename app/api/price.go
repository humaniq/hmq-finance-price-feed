package api

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
	"net/http"
	"strings"
)

const CtxSymbolKey = "symbol"
const CtxCurrencyKey = "currency"
const DefaultCurrency = "ETH"

func MayHaveStringListInQueryMiddlewareFunc(queryKey string, ctxKey string, caseCast int, delimiter string, defaults []string) func(next http.Handler) http.Handler {
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

func GetPricesFunc(state svc.PriceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		symbols := MustGetStringListFromCtx(ctx, CtxSymbolKey)
		currencies := MustGetStringListFromCtx(ctx, CtxCurrencyKey)
		if len(symbols) == 0 || len(currencies) == 0 {
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("invalid symbol/currency mapping"), http.StatusBadRequest)
			return
		}
		resultMap := make(map[string][]*svc.PriceRecord)
		for _, symbol := range symbols {
			prices := make([]*svc.PriceRecord, 0, len(currencies))
			for _, currency := range currencies {
				price, err := state.GetLatestSymbolPrice(ctx, symbol, currency)
				if err != nil {
					if !errors.Is(err, svc.ErrNoValue) {
						httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("internal error"), http.StatusInternalServerError)
						return
					}
					continue
				}
				prices = append(prices, price)
			}
			resultMap[symbol] = prices
		}
		httpext.JSON(w, httpapi.NewOkResponse().WithPayload(resultMap))
	}
}

func GetPriceForSymbolHandlerFunc(state svc.PriceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		symbol := httpapi.MustGetStringValueFromContext(ctx, CtxSymbolKey)
		currency, available := httpapi.GetStringValueFromContext(ctx, CtxCurrencyKey)
		if !available {
			currency = DefaultCurrency
		}
		price, err := state.GetLatestSymbolPrice(ctx, symbol, currency)
		if err != nil {
			if errors.Is(err, svc.ErrNoValue) {
				httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("not found"), http.StatusNotFound)
				return
			}
			httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("internal error"), http.StatusInternalServerError)
		}
		httpext.JSON(w, price)
	}
}

func GetSymbolPricesFunc(state svc.PriceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		symbol := httpapi.MustGetStringValueFromContext(ctx, CtxSymbolKey)
		var list []string
		if err := json.NewDecoder(r.Body).Decode(&list); err != nil || len(list) == 0 {
			list = defaultCurrencyList
		}
		result := make([]*svc.PriceRecord, 0, len(list))
		for _, currency := range list {
			price, err := state.GetLatestSymbolPrice(ctx, symbol, currency)
			if err != nil {
				if !errors.Is(err, svc.ErrNoValue) {
					httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("error reading price"), http.StatusInternalServerError)
					return
				}
				continue
			}
			result = append(result, price)
		}
		httpext.JSON(w, result)
	}
}

func GetPricesForListFunc(state svc.PriceGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		currency, available := httpapi.GetStringValueFromContext(ctx, CtxCurrencyKey)
		if !available {
			currency = DefaultCurrency
		}
		var list []string
		if err := json.NewDecoder(r.Body).Decode(&list); err != nil || len(list) == 0 {
			list = defaultSymbolList
		}
		result := make([]*svc.PriceRecord, 0, len(list))
		for _, symbol := range list {
			price, err := state.GetLatestSymbolPrice(ctx, symbol, currency)
			if err != nil {
				if !errors.Is(err, svc.ErrNoValue) {
					httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("error reading price"), http.StatusInternalServerError)
					return
				}
				continue
			}
			result = append(result, price)
		}
		httpext.JSON(w, result)
	}
}

var defaultSymbolList = []string{"ETH", "BTC", "USDT"}
var defaultCurrencyList = []string{"ETH", "USD", "EUR"}
