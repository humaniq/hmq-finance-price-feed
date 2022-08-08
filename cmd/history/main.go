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

	//cgSymbols = map[string]string{"usdt": "tether"}

	//log.Fatalf("%+v", cgSymbols)

	//cgSymbols = map[string]string{
	//	"usdt": "tether-avalanche-bridged-usdt-e",
	//}

	//smb := make(map[string]int)

	//assets, err := integrations.LoadPrices(ctx, currency)
	//if err != nil {
	//	log.Fatal(err)
	//}

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
		price0 := price.Value{
			Source:   "coingecko",
			Symbol:   key,
			Currency: currency,
		}
		//log.Printf("Price0: %+v\n", price0)
		records := make([]price.HistoryRecord, 0, len(ph.Prices))
		for index, a := range ph.Prices {
			ts := time.Unix(int64(a[0])/1000, int64(a[0])%1000)
			if index == 0 {
				price0.Price = a[1]
				price0.TimeStamp = ts
			}
			records = append(records, price.HistoryRecord{
				TimeStamp: ts,
				Price:     a[1],
			})
		}
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
				key: price.History{}.AddRecords(true, records...),
			},
		}
		log.Printf("ASSET: %+v", newAsset)
		log.Printf("HISTORY: %+v", newAsset.History)
		if err := backend.SavePrices(ctx, currency, &newAsset); err != nil {
			log.Fatal(err)
		}
	}
	log.Printf("NILS ARE: %+v", nilCurrencies)
}
