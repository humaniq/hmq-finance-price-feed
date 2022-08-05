package main

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"os"
	"strconv"
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

	if len(cfg.Providers.Coingeckos) > 0 {
		assets, err := config.AssetsFromFile(cfg.Coingecko.AssetsPath)
		if err != nil {
			app.Logger().Fatal(ctx, "FAIL GETTING COINGECKO ASSETS: %s", err)
			return
		}
		cg := prices.NewCoingecko(assets)
		for _, gecko := range cfg.Providers.Coingeckos {
			providerPool.AddProvider(prices.NewProvider(gecko.Name, cg.GetterFunc(gecko.Symbols, gecko.Currencies), gecko.Every()))
		}
	}

	//feed := make(chan []price.Value)

	storageWorker, err := prices.N

	consumer := prices.NewConsumer()
	consumer.AddWorker(&prices.LogWorker{})

	//for _, oracleCfg := range cfg.Consumers.PriceOracles {
	//	oracle, err := prices.NewPriceOracle(oracleCfg, cfg.EthNetworks.Networks)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
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
