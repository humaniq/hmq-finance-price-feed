package prices_old

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Messari struct {
	client *http.Client
}

func NewMessari() *Messari {
	return &Messari{}
}

func (m *Messari) httpClient() *http.Client {
	if m.client == nil {
		return http.DefaultClient
	}
	return m.client
}

func (m *Messari) GetPricesForSymbol(ctx context.Context, symbol string) (*MessariPricesValue, error) {
	url := fmt.Sprintf("https://data.messari.io/api/v1/assets/%s/metrics", strings.ToLower(symbol))
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	resp, err := m.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("status !ok")
	}
	var apiResp messariApiResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, err
	}
	return &MessariPricesValue{
		Eth: apiResp.Data.MarketData.PriceETH,
		Btc: apiResp.Data.MarketData.PriceBTC,
		Usd: apiResp.Data.MarketData.PriceUSD,
	}, nil
}

type messariApiResponse struct {
	Data messariApiResponseData `json:"data"`
}
type messariApiResponseData struct {
	Symbol     string                       `json:"symbol"`
	MarketData messariApiResponseMarketData `json:"market_data"`
}
type messariApiResponseMarketData struct {
	PriceUSD float64 `json:"price_usd"`
	PriceBTC float64 `json:"price_btc"`
	PriceETH float64 `json:"price_eth"`
}

type MessariPricesValue struct {
	Eth float64
	Btc float64
	Usd float64
}
