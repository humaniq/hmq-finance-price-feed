package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/humaniq/hmq-finance-price-feed/app/api"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		if logLevelNumeric, err := strconv.ParseUint(logLevel, 10, 8); err == nil {
			logger.InitDefault(uint8(logLevelNumeric))
		} else {
			logger.InitDefault(blogger.StringToLevel(logLevel))
		}
	}
	ctx := context.Background()

	priceCache, err := cache.NewLRU(1000)
	if err != nil {
		logger.Fatal(ctx, "priceCache init: %s", err)
		return
	}
	backend := svc.NewPriceStateSvc().WithCache(priceCache)

	if dsProjectId := os.Getenv("DATASTORE_PROJECT_ID"); dsProjectId != "" {
		ds, err := gds.NewClient(ctx, dsProjectId, "hmq_prices")
		if err != nil {
			logger.Fatal(ctx, "priceDS init: %s", err)
			return
		}
		backend = backend.WithGDSClient(ds)
	}
	back := svc.NewPriceStateEstimateWrapper(backend, "ETH")

	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Use(chim.Logger)
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins: []string{"https://*", "http://*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-auth-token"},
			MaxAge:         300,
		}))
		r.Route("/api/v1", func(r chi.Router) {
			if openapiPath := os.Getenv("OPENAPI_PATH"); openapiPath != "" {
				r.Group(func(r chi.Router) {
					//r.Use(cors.Handler(cors.Options{
					//	AllowedOrigins: []string{"https://*", "http://*"},
					//	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
					//	AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-auth-token"},
					//	MaxAge:         300,
					//}))
					r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
						http.ServeFile(w, r, openapiPath)
					})
				})
				r.Group(func(r chi.Router) {
					r.Use(
						httpapi.MustHaveStringValueInPathCtxMiddleware("symbol", api.CtxSymbolKey, httpapi.CaseToUpper),
						httpapi.MayHaveStringValueInQueryCtxMiddleware("currency", api.CtxCurrencyKey, httpapi.CaseToUpper),
					)
					r.Get("/price/{symbol}", api.GetPriceForSymbolHandlerFunc(back))
					r.Get("/prices/{symbol}", api.GetSymbolPricesFunc(back))
				})
				r.Group(func(r chi.Router) {
					r.Use(httpapi.MayHaveStringValueInQueryCtxMiddleware("currency", api.CtxCurrencyKey, httpapi.CaseToUpper))
					r.Get("/prices", api.GetPricesForListFunc(back))
				})
				r.Group(func(r chi.Router) {
					r.Use(api.MayHaveStringListInQueryMiddlewareFunc("symbol", api.CtxSymbolKey, httpapi.CaseToUpper, ",", []string{"ETH"}))
					r.Use(api.MayHaveStringListInQueryMiddlewareFunc("currency", api.CtxCurrencyKey, httpapi.CaseToUpper, ",", []string{"ETH", "USD", "EUR"}))
					r.Get("/prices/list", api.GetPricesFunc(back))
				})
			}
		})
	})

	sslHost := os.Getenv("APP_SSL_HOST")
	listenOn := os.Getenv("APP_LISTEN_ON")

	httpServer := &http.Server{Handler: router}

	if sslHost != "" {
		logger.Debug(ctx, "SSL enabled")
		if sslCacheDir := os.Getenv("SSL_CACHE_DIR"); sslCacheDir != "" {
			logger.Debug(ctx, "SSL cache")
			sslManager := &autocert.Manager{
				Prompt: autocert.AcceptTOS,
				HostPolicy: func(ctx context.Context, host string) error {
					if host == sslHost {
						return nil
					}
					return fmt.Errorf("acme/autocert: only %s host is allowed", sslHost)
				},
				Cache: autocert.DirCache(sslCacheDir),
			}
			httpServer.TLSConfig = &tls.Config{GetCertificate: sslManager.GetCertificate}
		}
		if sslDir := os.Getenv("SSL_DIR"); sslDir != "" {
			logger.Debug(ctx, "SSL existing certs")
			cert, err := tls.LoadX509KeyPair(filepath.Join(sslDir, sslHost, "fullchain.pem"), filepath.Join(sslDir, sslHost, "privkey.pem"))
			if err != nil {
				logger.Fatal(ctx, err.Error())
				return
			}
			httpServer.TLSConfig = &tls.Config{
				Certificates: []tls.Certificate{cert},
			}
		}

		if listenOn == "" {
			listenOn = ":443"
		}
		httpServer.Addr = listenOn

		logger.Info(ctx, "Prices API: listening (SSL=%s) on %s", sslHost, listenOn)

		if err := httpServer.ListenAndServeTLS("", ""); err != nil {
			logger.Fatal(ctx, err.Error())
			return
		}
		return
	}

	if listenOn == "" {
		listenOn = ":80"
	}

	httpServer.Addr = listenOn

	logger.Info(ctx, "Prices API: listening on %s", listenOn)

	if err := httpServer.ListenAndServe(); err != nil {
		logger.Fatal(ctx, err.Error())
	}

}
