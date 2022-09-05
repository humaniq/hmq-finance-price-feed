package main

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"github.com/humaniq/hmq-finance-price-feed/pkg/ethereum"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"os"
	"strconv"

	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
)

func main() {

	ctx := context.Background()

	logLevel := blogger.LevelDefault
	if logLevelEnv := os.Getenv("LOG_LEVEL"); logLevelEnv != "" {
		if logLevelNumeric, err := strconv.ParseUint(logLevelEnv, 10, 8); err == nil {
			logLevel = uint8(logLevelNumeric)
		} else {
			logLevel = blogger.StringToLevel(logLevelEnv)
		}
	}
	app.InitLogger(
		blogger.NewLog(
			[]blogger.LoggerMiddlewareFunc{
				blogger.LogLevelFilter(logLevel),
				blogger.LevelPrefix(),
				blogger.CurrentTimeFormat("(2006-01-02)(15:04:05MST)"),
				blogger.CtxStringValues("uid", "tag"),
			},
			blogger.NewIOWriterRouter(os.Stdout, os.Stderr, os.Stderr, true)),
	)

	environment := os.Getenv("ENVIRONMENT")
	if environment == "" {
		environment = "develop"
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = fmt.Sprintf("/etc/hmq/%s.price-feed.config.yaml", environment)
	}
	secretPath := os.Getenv("SECRET_PATH")
	if secretPath == "" {
		secretPath = "/secret"
	}

	cfg, err := config.PriceFeedConfigFromFile(configPath)
	if err != nil {
		app.Logger().Fatal(ctx, "ERROR READING CONFIG: %s", err)
		return
	}
	cfg.OverridesFromEnv()

	app.Logger().Info(ctx, "CONFIG: %+v", *cfg)

	providerPool := prices.NewProviderPool()

	if len(cfg.Providers) == 0 {
		app.Logger().Fatal(ctx, "NO PROVIDERS GIVEN")
		return
	}

	for _, providerCfg := range cfg.Providers {
		if providerCfg.GeoCurrency != nil {
			providerPool.AddProvider(
				prices.NewProvider(providerCfg.Name, prices.GeoCurrencyProviderFunc(providerCfg.GeoCurrency), providerCfg.Every()),
			)
		}
		if providerCfg.CoinGecko != nil {
			providerPool.AddProvider(
				prices.NewProvider(
					providerCfg.Name,
					prices.NewCoingecko(cfg.AssetsData.CoinGecko).GetterFunc(providerCfg.CoinGecko.Symbols, providerCfg.CoinGecko.Currencies),
					providerCfg.Every(),
				),
			)
		}
		if providerCfg.PancakeSwap != nil {
			providerPool.AddProvider(
				prices.NewProvider(
					providerCfg.Name,
					prices.NewPancakeSwap(providerCfg.PancakeSwap, cfg.AssetsData.EthNetworks["BSC"]).GetterFunc(),
					providerCfg.Every(),
				),
			)
		}
	}

	consumerState := prices.NewConsumerState(prices.SymbolCurrencyStateKey)

	consumer := prices.NewConsumer().WithFilters(
		prices.AnyOf(
			consumerState.TimeDeltaThresholdsFunc(cfg.Thresholds),
			consumerState.PercentThresholdsFunc(cfg.Thresholds),
		),
	)
	consumer.AddWorker(&prices.LogWorker{})

	for _, storageCfg := range cfg.Consumers {
		if storageCfg.GoogleDataStore != nil {
			ds, err := gds.NewClient(ctx, storageCfg.GoogleDataStore.ProjectID(), storageCfg.GoogleDataStore.PriceAssetsKind())
			if err != nil {
				app.Logger().Fatal(ctx, "FAIL GETTING Google Datastore Backend: %s", err)
				return
			}
			storageWorker := prices.NewStorageWriteWorker(storage.NewPricesDS(ds))
			consumer.AddWorker(storageWorker)
		}
		if storageCfg.PriceOracle != nil {
			network, found := cfg.AssetsData.EthNetworks[storageCfg.PriceOracle.NetworkKey]
			if !found {
				app.Logger().Fatal(ctx, "FAIL GETTING NETWORK %s", storageCfg.PriceOracle.NetworkKey)
				return
			}
			conn, err := ethereum.NewTransactConnection(network.RawUrl, network.ChainId, storageCfg.PriceOracle.ClientPrivateKey, 30000)
			if err != nil {
				app.Logger().Fatal(ctx, "FAIL ESTABLISHING ETH CONNECTION %s: %s", storageCfg.PriceOracle.NetworkKey, err)
				return
			}
			oracleWriter, err := ethereum.NewPriceOracleWriter(conn, storageCfg.PriceOracle.ContractAddressHex)
			if err != nil {
				app.Logger().Fatal(ctx, "FAIL CONNECTING CONTRACT %s", storageCfg.PriceOracle.ContractAddressHex)
				return
			}
			tokenValues := make([]price.Value, 0, len(storageCfg.Tokens))
			tokenMap := make(map[string]config.EthNetworkSymbolContract)
			for _, token := range storageCfg.Tokens {
				tokenValues = append(tokenValues, price.Value{
					Symbol:   token.Symbol,
					Currency: token.Currency,
				})

				tokenMap[token.Symbol] = network.Symbols[token.Symbol]
			}

			oracleState := prices.NewConsumerState(prices.SymbolCurrencyStateKey).WithValues(tokenValues...)
			oracle := prices.NewConsumerWorkerFilterWpapper(
				oracleState.ValueExists,
			).Wrap(prices.NewPriceOracleWriteWorker(tokenMap, oracleWriter))
			consumer.AddWorker(oracle)
		}
	}

	if err := consumer.Consume(ctx, providerPool.Feed()); err != nil {
		app.Logger().Fatal(ctx, "CONSUMER FAILED: %s", err)
		return
	}
	defer consumer.WaitForDone()

	if err := providerPool.Start(ctx); err != nil {
		app.Logger().Fatal(ctx, "PROVIDER START FAILED: %s", err)
		return
	}

}
