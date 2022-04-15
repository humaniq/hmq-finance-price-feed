package httpapi

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpext"
)

const (
	CaseSensitive = 0
	CaseToUpper   = 1
	CaseToLower   = 2
)

func MustHaveStringValueInPathCtxMiddleware(pathKey string, ctxKey string, caseCast int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := chi.URLParam(r, pathKey)
			if value == "" {
				httpext.AbortJSON(w, NewErrorResponse().WithPayload(fmt.Sprintf("no %s given", pathKey)), http.StatusBadRequest)
				return
			}
			value = toCaseGiven(value, caseCast)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey, value)))
		})
	}
}

func MustHaveStringValueInHeaderCtxMiddleware(headerKey string, ctxKey string, caseCast int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := r.Header.Get(headerKey)
			if value == "" {
				httpext.AbortJSON(w, NewErrorResponse().WithPayload(fmt.Sprintf("no %s given", headerKey)), http.StatusBadRequest)
				return
			}
			value = toCaseGiven(value, caseCast)
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey, value)))
		})
	}
}

func MayHaveStringValueInQueryCtxMiddleware(queryKey string, ctxKey string, caseCast int) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			value := r.URL.Query().Get(queryKey)
			if value != "" {
				ctx = context.WithValue(ctx, ctxKey, toCaseGiven(value, caseCast))
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func MustHaveStringCtxValueInListMiddleware(list []string, ctxKey string) func(next http.Handler) http.Handler {
	search := make(map[string]bool)
	for _, key := range list {
		search[key] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			value := MustGetStringValueFromContext(r.Context(), ctxKey)
			if !search[value] {
				httpext.AbortJSON(w, NewErrorResponse().WithPayload(fmt.Sprintf("%s not allowed", value)), http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func MustGetStringValueFromContext(ctx context.Context, ctxKey string) string {
	return ctx.Value(ctxKey).(string)
}
func GetStringValueFromContext(ctx context.Context, ctxKey string) (string, bool) {
	value := ctx.Value(ctxKey)
	if value == nil {
		return "", false
	}
	result, ok := value.(string)
	return result, ok
}

func toCaseGiven(value string, caseGiven int) string {
	switch caseGiven {
	case CaseToUpper:
		return strings.ToUpper(value)
	case CaseToLower:
		return strings.ToLower(value)
	default:
		return value
	}
}
