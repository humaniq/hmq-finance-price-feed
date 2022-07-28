package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type ChartValue struct {
	TimeStamp time.Time
	Value     float64
}

type MarketChart struct {
	MarketCaps []ChartValue
	Prices     []ChartValue
}

func GetMarketChart(ctx context.Context, client *http.Client, id string, currency string, days int) (*MarketChart, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/coins/%s/market_chart?vs_currency=%s&days=%d",
		id, currency, days)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrWrongStatus, resp.StatusCode)
	}
	defer resp.Body.Close()
	var body marketChartResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	chart := MarketChart{
		MarketCaps: make([]ChartValue, 0, len(body.MarketCaps)),
		Prices:     make([]ChartValue, 0, len(body.Prices)),
	}
	for _, value := range body.MarketCaps {
		chart.MarketCaps = append(chart.MarketCaps, ChartValue{
			TimeStamp: time.Unix(int64(value[0])/1000, int64(value[0])%1000),
			Value:     value[1],
		})
	}
	for _, value := range body.Prices {
		chart.Prices = append(chart.Prices, ChartValue{
			TimeStamp: time.Unix(int64(value[0])/1000, int64(value[0])%1000),
			Value:     value[1],
		})
	}
	return &chart, nil
}

type marketChartResponseBody struct {
	MarketCaps [][]float64 `json:"market_caps"`
	Prices     [][]float64 `json:"prices"`
}
