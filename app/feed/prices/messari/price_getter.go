package messari

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type PricesMessage struct {
	Eth float64
	Btc float64
	Usd float64
}

type PriceGetter struct {
	client *http.Client
}

func NewPriceGetter() *PriceGetter {
	return &PriceGetter{}
}

func (pg *PriceGetter) httpClient() *http.Client {
	if pg.client == nil {
		return http.DefaultClient
	}
	return pg.client
}

func (pg *PriceGetter) GetPricesForSymbol(ctx context.Context, symbol string) (*PricesMessage, error) {
	url := fmt.Sprintf("https://data.messari.io/api/v1/assets/%s/metrics", strings.ToLower(symbol))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := pg.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status !ok")
	}
	var apiResp apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	return &PricesMessage{
		Eth: apiResp.Data.MarketData.PriceETH,
		Btc: apiResp.Data.MarketData.PriceBTC,
		Usd: apiResp.Data.MarketData.PriceUSD,
	}, nil
}

type apiResponse struct {
	Data apiResponseData `json:"data"`
}
type apiResponseData struct {
	Symbol     string                `json:"symbol"`
	MarketData apiResponseMarketData `json:"market_data"`
}
type apiResponseMarketData struct {
	PriceUSD float64 `json:"price_usd"`
	PriceBTC float64 `json:"price_btc"`
	PriceETH float64 `json:"price_eth"`
}
