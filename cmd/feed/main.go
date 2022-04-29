package main

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/svc"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"github.com/humaniq/hmq-finance-price-feed/pkg/cache"
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

	priceCache, err := cache.NewLRU(1000)
	if err != nil {
		logger.Fatal(ctx, "priceCache init: %s", err)
		return
	}
	backend := svc.NewPriceStateSvc().WithCache(priceCache)
	logger.Info(ctx, "BACKEND INIT")

	if dsProjectId := os.Getenv("DATASTORE_PROJECT_ID"); dsProjectId != "" {
		ds, err := gds.NewClient(ctx, dsProjectId, "hmq_prices")
		if err != nil {
			logger.Fatal(ctx, "priceDS init: %s", err)
			return
		}
		backend = backend.WithGDSClient(ds)
		logger.Info(ctx, "DS storage added")
	}

	chainId, err := strconv.Atoi(os.Getenv("CONTRACT_PRICES_CHAIN_ID"))
	if err != nil {
		logger.Fatal(ctx, "chainID fail: %s", err.Error())
		return
	}

	deltas := make(map[string]int)
	deltas["ETH"] = 1

	contractSetter, err := svc.NewContractPriceSetter(
		os.Getenv("CONTRACT_PRICES_URL"),
		int64(chainId),
		os.Getenv("CONTRACT_PRICES_ADDRESS"),
		os.Getenv("CONTRACT_PRICES_PRIVATE_KEY"),
	)
	pricesContractConsumer := svc.NewContractPricesConsumer().
		WithGetter(backend).
		WithSetter(contractSetter).
		WithFilters(svc.FilterDeltaFunc(backend, deltas, time.Hour))

	pricesStorageConsumer := svc.NewStoreConsumer(backend).
		WithNext(pricesContractConsumer).
		WithFilters(svc.FilterDeltaFunc(backend, deltas, time.Hour))

	if err := pricesContractConsumer.Start(); err != nil {
		logger.Fatal(ctx, err.Error())
		return
	}
	defer pricesContractConsumer.WaitForDone()
	if err := pricesStorageConsumer.Start(); err != nil {
		logger.Fatal(ctx, err.Error())
		return
	}
	defer pricesStorageConsumer.WaitForDone()

	messariTickerDuration := time.Minute * 5
	if tickerSeconds, err := strconv.Atoi(os.Getenv("MESSARI_TICKER_SECONDS")); err == nil {
		messariTickerDuration = time.Second * time.Duration(tickerSeconds)
	}

	messariTokenList := strings.Split(os.Getenv("MESSARI_TOKEN_LIST"), ",")
	if len(messariTokenList) == 0 {
		messariTokenList = []string{"ETH", "USDT", "BTC"}
	}

	messariPricesProvider := svc.NewMessariPriceProvider(messariTickerDuration, messariTokenList)

	go messariPricesProvider.Provide(ctx, pricesContractConsumer.Lease())

	coingeckoClient, err := prices.CoinGeckoFromFile(os.Getenv("COINGECKO_CONFIG_FILE"))
	if err != nil {
		logger.Fatal(ctx, err.Error())
		return
	}
	coingeckoPricesProvider := svc.NewCoinGeckoProvider(time.Minute, coingeckoClient, strings.Split(os.Getenv("COINGECKO_SYMBOL_LIST"), ","), strings.Split("COINGECKO_CURRENCY_LIST", ","))
	coingeckoPricesProvider.Provide(ctx, pricesContractConsumer.Lease())

}
