package main

import (
	"context"
	"errors"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/feed"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/app/state"
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

	pricesState := make(map[string]*state.Prices)
	for _, currency := range cfg.Assets {
		currencyState, err := backend.LoadPrices(ctx, currency)
		if err != nil {
			if !errors.Is(err, storage.ErrNotFound) {
				logger.Fatal(ctx, "prices state init: %s", err)
				return
			}
			currencyState = state.NewPrices(currency)
		}
		currencyState = currencyState.WithCommitFilters(
			state.CommitTimestampFilterFunc(),
			state.CommitCurrenciesFilterFunc(map[string]bool{currency: true}),
			state.CommitPricePercentDiffFilterFinc(cfg.Diffs),
		)
		pricesState[currency] = currencyState
	}

	dsConsumer := feed.NewStorageConsumer("DS", backend, pricesState)

	coingeckoClient, err := prices.CoinGeckoFromFile(os.Getenv("COINGECKO_CONFIG_FILE"))
	if err != nil {
		logger.Fatal(ctx, err.Error())
		return
	}

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

		contractConsumer := feed.NewStorageConsumer("CONTRACT", contractBackend, pricesState)
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
				if err := feed.NewCoinGeckoProvider(time.Minute*5, coingeckoClient, provider.Symbols, provider.Currencies).
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
