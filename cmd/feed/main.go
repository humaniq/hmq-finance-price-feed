package main

import (
	"context"
	"fmt"
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
	}

	//if len(cfg.Providers.Coingeckos) > 0 {
	//	assets, err := config.AssetsFromFile(cfg.Coingecko.AssetsPath)
	//	if err != nil {
	//		app.Logger().Fatal(ctx, "FAIL GETTING COINGECKO ASSETS: %s", err)
	//		return
	//	}
	//	cg := prices.NewCoingecko(assets)
	//	for _, gecko := range cfg.Providers.Coingeckos {
	//		providerPool.AddProvider(prices.NewProvider(gecko.Name, cg.GetterFunc(gecko.Symbols, gecko.Currencies), gecko.Every()))
	//	}
	//}

	//feed := make(chan []price.Value)

	consumer := prices.NewConsumer()
	consumer.AddWorker(&prices.LogWorker{})

	//for _, storageCfg := range cfg.Consumers.StorageConsumers {
	//	ds, err := gds.NewClient(ctx, storageCfg.Datastore.ProjectID(), storageCfg.Datastore.PriceAssetsKind())
	//	if err != nil {
	//		app.Logger().Fatal(ctx, "FAIL GETTING Google Datastore Backend")
	//		return
	//	}
	//	consumerState := prices.NewConsumerState(prices.SymbolCurrencyStateKey)
	//	storageWorker := prices.NewConsumerWorkerFilterWpapper(prices.AnyOf(
	//		consumerState.TimeDeltaFunc(storageCfg.TimeDelta()),
	//		consumerState.PercentThresholdFunc(storageCfg.Thresholds.Symbols, storageCfg.Thresholds.Default, prices.SymbolStateKey),
	//	)).Wrap(prices.NewStorageWriteWorker(storage.NewPricesDS(ds)))
	//	consumer.AddWorker(storageWorker)
	//}
	//for _, oracleCfg := range cfg.Consumers.PriceOracles {
	//	network, found := cfg.EthNetworks.Networks[oracleCfg.NetworkUid]
	//	if !found {
	//		app.Logger().Fatal(ctx, "FAIL GETTING NETWORK %s", oracleCfg.NetworkUid)
	//		return
	//	}
	//	conn, err := ethereum.NewTransactConnection(network.RawUrl, network.ChainId, oracleCfg.ClientPrivateKey, 30000)
	//	if err != nil {
	//		app.Logger().Fatal(ctx, "FAIL ESTABLISHING ETH CONNECTION %s: %s", oracleCfg.NetworkUid, err)
	//		return
	//	}
	//	oracleWriter, err := ethereum.NewPriceOracleWriter(conn, oracleCfg.ContractAddressHex)
	//	if err != nil {
	//		app.Logger().Fatal(ctx, "FAIL CONNECTING CONTRACT %s", oracleCfg.ContractAddressHex)
	//		return
	//	}
	//	tokenMap := make(map[string]config.EthNetworkSymbolContract)
	//	thresholdsMap := make(map[string]float64)
	//	tokenValues := make([]price.Value, 0, len(oracleCfg.Tokens))
	//	for _, token := range oracleCfg.Tokens {
	//		tokenCfg, found := network.Symbols[token.Symbol]
	//		if !found {
	//			app.Logger().Fatal(ctx, "FAIL FINDING TOKEN ADDRESS for %s in %s", token.Symbol, network.Name)
	//			return
	//		}
	//		if token.PercentThreshold > 0 {
	//			thresholdsMap[fmt.Sprintf("%s-%s", token.Symbol, token.Currency)] = token.PercentThreshold
	//		} else {
	//			app.Logger().Fatal(ctx, "PERCENT THRESHOLD IS 0 for (%s)%s-%s", network.Name, token.Symbol, token.Currency)
	//			return
	//		}
	//		tokenMap[token.Symbol] = config.EthNetworkSymbolContract{
	//			AddressHex: tokenCfg.AddressHex,
	//			Decimals:   tokenCfg.Decimals,
	//		}
	//		tokenValues = append(tokenValues, price.Value{
	//			Symbol:   token.Symbol,
	//			Currency: token.Currency,
	//		})
	//	}
	//	consumerState := prices.NewConsumerState(prices.SymbolCurrencyStateKey).WithValues(tokenValues...)
	//	oracle := prices.NewConsumerWorkerFilterWpapper(prices.AllOf(
	//		consumerState.ValueExists,
	//		consumerState.PercentThresholdFunc(thresholdsMap, 0, prices.SymbolCurrencyStateKey),
	//	)).Wrap(prices.NewPriceOracleWriteWorker(tokenMap, oracleWriter))
	//	consumer.AddWorker(oracle)
	//}

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
