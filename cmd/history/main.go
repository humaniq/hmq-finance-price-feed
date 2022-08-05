package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/app/storage"
	"github.com/humaniq/hmq-finance-price-feed/pkg/gds"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"sort"
	"time"
)

type PricesHistory struct {
	Prices [][]float64 `json:"prices"`
}

func main() {

	ctx := context.Background()

	gdsClient, err := gds.NewClient(ctx, "humaniq-168420", "hmq_price_assets")
	if err != nil {
		log.Fatal(err)
	}
	backend := storage.NewPricesDS(gdsClient)

	currency := os.Getenv("CURRENCY")

	cfg := struct {
		Symbols map[string]string `yaml:"symbols"`
	}{
		Symbols: make(map[string]string),
	}

	file, err := os.Open("etc/coingecko.assets.yaml")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		log.Fatal(err)
	}
	cgSymbols := cfg.Symbols

	//log.Fatalf("%+v", cgSymbols)

	//cgSymbols = map[string]string{
	//	"usdt": "tether-avalanche-bridged-usdt-e",
	//}

	//smb := make(map[string]int)

	//assets, err := integrations.LoadPrices(ctx, currency)
	//if err != nil {
	//	log.Fatal(err)
	//}

	dayTimeLimit := time.Now().Add(time.Hour * (-24))
	weekTimeLimit := dayTimeLimit.Add(time.Hour * (-24) * 7)
	monthTimeLimit := weekTimeLimit.Add(time.Hour * (-24) * 31)

	log.Println(dayTimeLimit)

	var nilCurrencies []price.Value

	counter := 0

	for key, val := range cgSymbols {
		if counter == 45 {
			counter = 0
			<-time.Tick(time.Minute)
		}
		counter++
		log.Println(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=40", val, currency))
		h, err := http.Get(fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=40", val, currency))
		if err != nil {
			log.Fatal(err)
		}
		var ph PricesHistory
		if err := json.NewDecoder(h.Body).Decode(&ph); err != nil {
			log.Fatal(err)
		}
		h.Body.Close()
		//log.Printf("%+v", ph)
		sort.Slice(ph.Prices, func(i, j int) bool {
			return ph.Prices[i][0] > ph.Prices[j][0]
		})
		var historyRecords []price.HistoryRecord
		var latestHistoryRecord price.HistoryRecord
		delta := -time.Hour
		price0 := price.Value{
			Source:   "coingecko",
			Symbol:   key,
			Currency: currency,
		}
		//log.Printf("Price0: %+v\n", price0)
		for index, a := range ph.Prices {
			ts := time.Unix(int64(a[0])/1000, int64(a[0])%1000)
			if index == 0 {
				price0.Price = a[1]
				price0.TimeStamp = ts
			}
			if ts.Before(dayTimeLimit) {
				delta = -time.Hour * 4
			}
			if ts.Before(weekTimeLimit) {
				delta = -time.Hour * 6
			}
			if ts.Before(monthTimeLimit) {
				break
			}
			if latestHistoryRecord.TimeStamp.IsZero() ||
				ts.Before(latestHistoryRecord.TimeStamp.Add(delta)) {
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
		log.Printf("historyRecords: %d", len(historyRecords))
		log.Printf("Price1: %+v\n", price0)
		if price0.Price == 0 {
			nilCurrencies = append(nilCurrencies, price0)
			continue
		}
		newAsset := price.Asset{
			Name: currency,
			Prices: map[string]price.Value{
				key: price0,
			},
			History: map[string]price.History{
				key: historyRecords,
			},
		}
		if err := backend.SavePrices(ctx, currency, &newAsset); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("NILS ARE: %+v", nilCurrencies)

	//universe-token

	//tt := 0
	//for _, asset := range cfg.Assets {
	//	assets, err := integrations.LoadPrices(ctx, asset)
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
	//						Value:     a[1],
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
	//			if err := integrations.SavePrices(ctx, asset, &newAsset); err != nil {
	//				log.Fatal(err)
	//			}
	//			//log.Fatal("DONE")
	//		}
	//	}
	//}

	//log.Printf("%+v", smb)

}
