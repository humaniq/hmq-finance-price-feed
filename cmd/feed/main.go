package main

import (
	"context"
	"errors"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/state"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/feed"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
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

	configPath := os.Getenv("CONFIG_FILE_PATH")
	if configPath == "" {
		configPath = "/etc/hmq/price-feed.yaml"
	}
	cfg, err := config.FeedConfigFromFile(configPath)
	if err != nil {
		logger.Fatal(ctx, "error getting config: %s", err)
		return
	}

	dsKind := os.Getenv("DATASTORE_PRICES_KIND")
	if dsKind == "" {
		dsKind = "hmq_current_prices"
	}
	gdsClient, err := gds.NewClient(ctx, os.Getenv("DATASTORE_PROJECT_ID"), dsKind)
	if err != nil {
		logger.Fatal(ctx, "gdsClient init: %s", err)
		return
	}
	backend := storage.NewPricesDS(gdsClient)

	pricesState := make(map[string]*state.AssetCommitter)
	for _, currency := range cfg.Assets {
		currencyState, err := backend.LoadPrices(ctx, currency)
		if err != nil {
			if !errors.Is(err, app.ErrNotFound) {
				logger.Fatal(ctx, "prices state init: %s", err)
				return
			}
			currencyState = price.NewAsset(currency)
		}
		pricesState[currency] = state.NewAssetCommitter(currencyState).WithFilters(
			state.CommitValueCurrenciesFilterFunc(map[string]bool{currency: true}),
			state.CommitValuePricePercentDiffFilterFinc(cfg.Diffs),
		)
	}

	dsConsumer := feed.NewStorageConsumer("DS", backend, pricesState)

	if contractUrl := os.Getenv("CONTRACT_PRICES_URL"); contractUrl != "" {
		chainIdString := os.Getenv("CONTRACT_CHAIN_ID")
		chainId, err := strconv.ParseInt(chainIdString, 10, 64)
		if err != nil {
			chainId = 1337
		}
		contractBackend, err := storage.NewPricesContractSetter(
			contractUrl,
			chainId,
			os.Getenv("CONTRACT_PRICES_ADDRESS"),
			os.Getenv("CONTRACT_PRICES_PRIVATE_KEY"),
		)
		if err != nil {
			logger.Fatal(ctx, "err getting contract backend: %s", err)
			return
		}

		ps := make(map[string]*state.AssetCommitter)
		ps["usd"] = pricesState["usd"]
		ps["eur"] = pricesState["eur"]

		contractConsumer := feed.NewStorageConsumer("CONTRACT", contractBackend, ps)
		go contractConsumer.Run()
		defer contractConsumer.WaitForDone()
		dsConsumer = dsConsumer.WithNext(
			contractConsumer,
		)
	}

	go dsConsumer.Run()
	defer dsConsumer.WaitForDone()

	var wg sync.WaitGroup
	wg.Add(len(cfg.Providers))
	defer wg.Wait()

	for _, provider := range cfg.Providers {
		switch provider.Type {
		case "coingecko":
			go func(name string) {
				defer wg.Done()
				defer dsConsumer.Release()
				if err := feed.NewCoinGeckoProvider(time.Minute*5, prices.NewCoinGecko(), provider.Symbols, provider.Currencies).
					Provide(ctx, dsConsumer.Lease()); err != nil {
					logger.Fatal(ctx, "%s provider fail: %s", name, err.Error())
				}
			}(provider.Name)
		case "geocurrency":
			go func(name string) {
				defer wg.Done()
				defer dsConsumer.Release()
				geoCurrencyProvider := feed.NewGeoCurrencyPriceProvider(time.Minute*5, prices.NewIPCurrencyAPI(os.Getenv("GEO_CURRENCY_KEY")), provider.Symbols, provider.Currencies)
				if err := geoCurrencyProvider.Provide(ctx, dsConsumer.Lease()); err != nil {
					logger.Fatal(ctx, "%s provider fail: %s", name, err.Error())
				}
			}(provider.Name)
		default:
			logger.Fatal(ctx, "UNKNOWN PROVIDER: %s", provider.Type)
			return
		}
	}

}
