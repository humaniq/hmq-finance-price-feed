package httpapi

import (
	"context"
	"github.com/google/uuid"
	"net/http"
)

const CtxRequestUidKey = "requestUid"

func RequestUidMiddleware(ctxKey string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey, uuid.New().String())))
		})
	}
}
