package api

import (
	"encoding/json"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
	"net/http"
)

const CtxSymbolKey = "symbol"
const CtxCurrencyKey = "currency"
const DefaultCurrency = "ETH"

func GetPriceForSymbolHandlerFunc(state *svc.PriceSvc) http.HandlerFunc {
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

func GetPricesForListFunc(state *svc.PriceSvc) http.HandlerFunc {
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
			if err != nil && !errors.Is(err, svc.ErrNoValue) {
				httpext.AbortJSON(w, httpapi.NewErrorResponse().WithPayload("error reading price"), http.StatusInternalServerError)
				return
			}
			result = append(result, price)
		}
		httpext.JSON(w, result)
	}
}

var defaultSymbolList = []string{"ETH", "BTC", "USDT"}
