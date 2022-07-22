package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"log"
	"net/http"
	"sort"
	"time"
)

type PricesHistory struct {
	Prices [][]float64 `json:"prices"`
}

func main() {

	ctx := context.Background()

	cfg, err := config.FeedConfigFromFile("etc/price-feed.config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	gdsClient, err := gds.NewClient(ctx, "humaniq-168420", "hmq_price_assets")
	if err != nil {
		log.Fatal(err)
	}
	backend := storage.NewPricesDSv2(gdsClient)

	currency := "usd"

	cgSymbols := make(map[string]string)
	for _, v := range cfg.Providers {
		if v.Name != "coingecko" {
			continue
		}
		for key, val := range v.Symbols {
			cgSymbols[val] = key
		}
	}

	//cgSymbols = map[string]string{
	//	"usdt": "tether-avalanche-bridged-usdt-e",
	//}

	//smb := make(map[string]int)

	assets, err := backend.LoadPrices(ctx, currency)
	if err != nil {
		log.Fatal(err)
	}

	historyTimeLimit := time.Now().Add(time.Hour * 24 * (-40))
	log.Println(historyTimeLimit)

	for key, val := range cgSymbols {
		log.Println(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=10", val, currency))
		h, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=10", val, currency))
		if err != nil {
			log.Fatal(err)
		}
		var ph PricesHistory
		if err := json.NewDecoder(h.Body).Decode(&ph); err != nil {
			log.Fatal(err)
		}
		h.Body.Close()
		sort.Slice(ph.Prices, func(i, j int) bool {
			return ph.Prices[i][0] > ph.Prices[j][0]
		})
		var historyRecords []price.HistoryRecord
		var latestHistoryRecord price.HistoryRecord
		price0 := assets.Prices[key]
		log.Printf("Price0: %+v\n", price0)
		for index, a := range ph.Prices {
			ts := time.Unix(int64(a[0])/1000, int64(a[0])%1000)
			log.Println(fmt.Sprintf("%f", a[0]))
			log.Println(ts)
			log.Println(a[1])
			if index == 0 {
				price0.Price = a[1]
				price0.TimeStamp = ts
			}
			if ts.Before(historyTimeLimit) {
				break
			}
			if latestHistoryRecord.TimeStamp.IsZero() ||
				ts.Before(latestHistoryRecord.TimeStamp.Add(-time.Hour*4)) {
				latestHistoryRecord = price.HistoryRecord{
					TimeStamp: ts,
					Price:     a[1],
				}
				historyRecords = append(historyRecords, latestHistoryRecord)
			}
		}
		sort.Slice(historyRecords, func(i, j int) bool {
			return historyRecords[i].TimeStamp.Before(historyRecords[j].TimeStamp)
		})
		log.Printf("history_records(%s:%s): %+v\n", key, val, historyRecords)
		log.Printf("Price1: %+v\n", price0)
		newAsset := price.Asset{
			Name: currency,
			Prices: map[string]price.Value{
				key: price0,
			},
			History: map[string]price.History{
				key: historyRecords,
			},
		}
		log.Printf("NEW_ASSET: %+v\n", newAsset)
		if err := backend.SavePrices(ctx, currency, &newAsset); err != nil {
			log.Fatal(err)
		}
	}

	//universe-token

	//tt := 0
	//for _, asset := range cfg.Assets {
	//	assets, err := backend.LoadPrices(ctx, asset)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	for symbol, history := range assets.History {
	//		counter := 0
	//		for _, hval := range history {
	//			if hval.TimeStamp.After(time.Date(2022, 6, 1, 0, 0, 0, 0, time.UTC)) {
	//				counter++
	//			}
	//		}
	//		if counter < 10 {
	//			smb[symbol] = smb[symbol] + 1
	//			log.Printf("%s=>%s", asset, symbol)
	//			log.Println(cgSymbols[symbol])
	//			if symbol == "stbz" {
	//				continue
	//			}
	//			//https://api.coingecko.com/api/v3/coins/usd-coin/market_chart?vs_currency=usd&days=10
	//			h, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=10", cgSymbols[symbol], asset))
	//			if err != nil {
	//				log.Fatal(err)
	//			}
	//			var ph PricesHistory
	//			if err := json.NewDecoder(h.Body).Decode(&ph); err != nil {
	//				log.Fatal(err)
	//			}
	//			h.Body.Close()
	//			var historyRecords []price.HistoryRecord
	//			latestHistoryRecord := price.HistoryRecord{}
	//			for _, a := range ph.Prices {
	//				ts := time.Unix(int64(a[0])/1000, int64(a[0])%1000)
	//				log.Printf("%s - %f\n", ts, a[1])
	//				if ts.After(latestHistoryRecord.TimeStamp.Add(time.Hour * 24)) {
	//					latestHistoryRecord = price.HistoryRecord{
	//						TimeStamp: ts,
	//						Price:     a[1],
	//					}
	//					historyRecords = append(historyRecords, latestHistoryRecord)
	//				}
	//			}
	//			log.Printf("history_records: %+v\n", historyRecords)
	//			price0 := assets.Prices[symbol]
	//			newAsset := price.Asset{
	//				Name: asset,
	//				Prices: map[string]price.Value{
	//					symbol: price0,
	//				},
	//				History: map[string]price.History{
	//					symbol: historyRecords,
	//				},
	//			}
	//			tt++
	//			if tt == 10 {
	//				tt = 0
	//				<-time.Tick(time.Minute)
	//			}
	//			if err := backend.SavePrices(ctx, asset, &newAsset); err != nil {
	//				log.Fatal(err)
	//			}
	//			//log.Fatal("DONE")
	//		}
	//	}
	//}

	//log.Printf("%+v", smb)

}
