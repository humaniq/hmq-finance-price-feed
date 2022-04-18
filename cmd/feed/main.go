package main

import (
	"context"
	"log"
	"os"
	"strconv"
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
	backend := svc.NewPriceSvc().WithCache(priceCache)

	if dsProjectId := os.Getenv("DATASTORE_PROJECT_ID"); dsProjectId != "" {
		ds, err := gds.NewClient(ctx, dsProjectId, "hmq_prices")
		if err != nil {
			logger.Fatal(ctx, "priceDS init: %s", err)
			return
		}
		backend = backend.WithGDSClient(ds)
	}

	feed := svc.NewPriceFeed().
		WithConsumerChan("log", svc.ConsumerForLog()).
		WithConsumerChan("ds", svc.ConsumerForDS(backend))
	if err := feed.Start(); err != nil {
		log.Fatal(err)
	}
	defer feed.WaitForDone()

	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			feed.Queue() <- svc.PriceRecord{
				Source:    "test",
				Symbol:    "ETH",
				Currency:  "USD",
				Price:     3400,
				TimeStamp: time.Now(),
			}
		}
	}()

}
