package main

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
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

	nativeCurrencyList := strings.Split(os.Getenv("NATIVE_CURRENCY_LIST"), ",")
	if len(nativeCurrencyList) == 0 || (len(nativeCurrencyList) == 1 && nativeCurrencyList[0] == "") {
		nativeCurrencyList = config.DefaultNativeCurrencyList
	}

	pricesState := make(map[string]*state.Prices)
	for _, currency := range nativeCurrencyList {
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
			state.CommitSymbolsFilterFunc(config.KnownSymbolsChecker()),
			state.CommitPricePercentDiffFilterFinc(config.DefaultSymbolDiffs()),
		)
		pricesState[currency] = currencyState
	}

	dsConsumer := feed.NewStorageConsumer("DS", backend, pricesState)

	coingeckoClient, err := prices.CoinGeckoFromFile(os.Getenv("COINGECKO_CONFIG_FILE"))
	if err != nil {
		logger.Fatal(ctx, err.Error())
		return
	}
	logger.Trace(ctx, "%+v", coingeckoClient)

	smb := strings.Split(os.Getenv("COINGECKO_SYMBOL_LIST"), ",")
	for index, value := range smb {
		smb[index] = strings.ToLower(value)
	}
	cur := strings.Split(os.Getenv("COINGECKO_CURRENCY_LIST"), ",")
	for index, value := range cur {
		cur[index] = strings.ToLower(value)
	}

	go dsConsumer.Run()
	defer dsConsumer.WaitForDone()

	coingeckoProvider := feed.NewCoinGeckoProvider(time.Minute*5, coingeckoClient, smb, cur)
	if err := coingeckoProvider.Provide(ctx, dsConsumer.Lease()); err != nil {
		logger.Fatal(ctx, "provider fail: %s", err)
		dsConsumer.Release()
		return
	}

}
