package main

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/prices"
	"log"
)

func main() {

	assets := config.Assets{
		Currencies: map[string]string{"usd": "usd"},
		Symbols:    map[string]string{"wbgl": "bitgesell", "busd": "busd"},
	}

	cg := prices.NewCoingecko(&assets)
	fn := cg.GetterFunc([]string{"wbgl", "busd"}, []string{"usd"})

	prices, err := fn(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PRICES: %+v\n", prices)
	for _, price := range prices {
		log.Printf("%+v\n", price)
	}

}
