package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	chim "github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/humaniq/hmq-finance-price-feed/app/api"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"github.com/humaniq/hmq-finance-price-feed/pkg/httpapi"
	"golang.org/x/crypto/acme/autocert"
)

func main() {

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		if logLevelNumeric, err := strconv.ParseUint(logLevel, 10, 8); err == nil {
			app.InitDefaultLogger(uint8(logLevelNumeric))
		} else {
			app.InitDefaultLogger(blogger.StringToLevel(logLevel))
		}
	}
	ctx := context.Background()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "/etc/hmq/price-api.config.yaml"
	}
	cfg, err := config.ApiConfigFromFile(configPath)
	if err != nil {
		app.Logger().Fatal(ctx, "error getting config: %s", err)
		return
	}

	app.Logger().Info(ctx, "STARTING WITH: %+v", cfg)

	gdsClient, err := gds.NewClient(ctx, cfg.Backend.GoogleDataStore.ProjectID(), cfg.Backend.GoogleDataStore.PriceAssetsKind())
	if err != nil {
		app.Logger().Fatal(ctx, "gdsClient init: %s", err)
		return
	}
	dsBackend := storage.NewPricesDS(gdsClient)

	backend := storage.NewInMemory(time.Minute*10).Wrap(dsBackend).Warm(context.Background(), cfg.Assets, time.Minute*9)

	router := chi.NewRouter()
	router.Use(httpapi.RequestUidMiddleware(httpapi.CtxRequestUidKey))
	router.Group(func(r chi.Router) {
		r.Use(
			chim.Logger,
			cors.Handler(cors.Options{
				AllowedOrigins: []string{"https://*", "http://*"},
				AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
				AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-auth-token"},
				MaxAge:         300,
			}),
		)
		r.Route("/api/v1", func(r chi.Router) {
			if openapiPath := os.Getenv("OPENAPI_PATH"); openapiPath != "" {
				r.Group(func(r chi.Router) {
					r.Get("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
						http.ServeFile(w, r, openapiPath)
					})
				})
			}
			r.Group(func(r chi.Router) {
				r.Use(api.MustHaveStringListInQueryOrDefaultsMiddlewareFunc("symbol", api.CtxSymbolKey, httpapi.CaseToLower, ",", []string{"eth"}))
				r.Use(api.MustHaveStringListInQueryOrDefaultsMiddlewareFunc("currency", api.CtxCurrencyKey, httpapi.CaseToLower, ",", []string{"eth", "usd", "eur", "rub"}))
				r.Get("/prices/list", api.GetPricesFunc(svc.NewPrices(backend).WithMapping(cfg.Currencies)))
			})
		})
	})

	sslHost := os.Getenv("APP_SSL_HOST")
	listenOn := cfg.API.Listen()

	httpServer := &http.Server{Handler: router}

	if sslHost != "" {
		app.Logger().Debug(ctx, "SSL enabled")
		if sslCacheDir := os.Getenv("SSL_CACHE_DIR"); sslCacheDir != "" {
			app.Logger().Debug(ctx, "SSL cache")
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
			app.Logger().Debug(ctx, "SSL existing certs")
			cert, err := tls.LoadX509KeyPair(filepath.Join(sslDir, sslHost, "fullchain.pem"), filepath.Join(sslDir, sslHost, "privkey.pem"))
			if err != nil {
				app.Logger().Fatal(ctx, err.Error())
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

		app.Logger().Info(ctx, "Prices Config: listening (SSL=%s) on %s", sslHost, listenOn)

		if err := httpServer.ListenAndServeTLS("", ""); err != nil {
			app.Logger().Fatal(ctx, err.Error())
			return
		}
		return
	}

	if listenOn == "" {
		listenOn = ":80"
	}

	httpServer.Addr = listenOn

	app.Logger().Info(ctx, "Prices Config: listening on %s", listenOn)

	if err := httpServer.ListenAndServe(); err != nil {
		app.Logger().Fatal(ctx, err.Error())
	}

}
