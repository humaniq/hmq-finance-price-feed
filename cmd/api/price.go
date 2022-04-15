package api

import (
	"errors"
	"net/http"

	"github.com/humaniq/hmq-finance-price-feed/cmd/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
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
