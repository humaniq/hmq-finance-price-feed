package main

import (
	"context"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/blogger"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
)

func main() {

	ctx := context.Background()

	if logLevel := os.Getenv("LOG_LEVEL"); logLevel != "" {
		if logLevelNumeric, err := strconv.ParseUint(logLevel, 10, 8); err == nil {
			logger.InitDefault(uint8(logLevelNumeric))
		} else {
			logger.InitDefault(blogger.StringToLevel(logLevel))
		}
	} else {
		logger.InitDefault(blogger.LevelInfo)
	}

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
		logger.Fatal(ctx, "ERROR READING CONFIG: %s", err)
		return
	}
	cfg.OverridesFromEnv()

	logger.Info(ctx, "CONFIG: %+v", cfg)

	var providers []*prices.Provider
	if len(cfg.Providers.Coingeckos) > 0 {
		assets, err := config.AssetsFromFile(cfg.Coingecko.AssetsPath)
		if err != nil {
			logger.Fatal(ctx, "FAIL GETTING COINGECKO ASSETS: %s", err)
			return
		}
		cg := prices.NewCoingecko(assets)
		for _, gecko := range cfg.Providers.Coingeckos {
			providers = append(providers, prices.NewProvider(gecko.Name, cg.GetterFunc(gecko.Symbols, gecko.Currencies), gecko.Every()))
		}
	}

	feed := make(chan []price.Value)

	consumer := prices.NewConsumer()

	for _, oracleCfg := range cfg.Consumers.PriceOracles {
		oracle, err := prices.NewPriceOracle(oracleCfg, cfg.EthNetworks.Networks)
		if err != nil {
			log.Fatal(err)
		}
		consumer.AddWorker(oracle)
	}

	if err := consumer.Consume(ctx, feed); err != nil {
		logger.Fatal(ctx, "CONSUMER FAILED: %s", err)
		return
	}
	defer consumer.WaitForDone()

	var wg sync.WaitGroup
	wg.Add(len(providers))
	for _, provider := range providers {
		go func(p *prices.Provider) {
			defer wg.Done()
			if err := p.Provide(ctx, feed); err != nil {
				logger.Fatal(ctx, "PROVIDER %s START ERROR: %s", p.Name(), err)
			}
			p.WaitForDone()
		}(provider)
	}
	wg.Wait()

}
