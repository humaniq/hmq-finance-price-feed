package prices

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var ErrWrongResponseStatusCode = errors.New("wrong response status code")

type IPCurrencyAPI struct {
	key        string
	httpClient *http.Client
}

func NewIPCurrencyAPI(key string) *IPCurrencyAPI {
	return &IPCurrencyAPI{key: key}
}
func (ip *IPCurrencyAPI) client() *http.Client {
	if ip.httpClient == nil {
		return http.DefaultClient
	}
	return ip.httpClient
}
func (ip *IPCurrencyAPI) WithHttpClient(client *http.Client) *IPCurrencyAPI {
	ip.httpClient = client
	return ip
}
func (ip *IPCurrencyAPI) GetConversionRates(ctx context.Context, baseCurrency string, amount float64, currencies ...string) (IPCurrencyRates, error) {
	url := fmt.Sprintf("https://api.getgeoapi.com/v2/currency/convert?api_key=%s&from=%s&amount=%f&format=json", ip.key, baseCurrency, amount)
	if len(currencies) > 0 {
		url = fmt.Sprintf("%s&to=%s", url, strings.Join(currencies, ","))
	}
	resp, err := ip.client().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, ErrWrongResponseStatusCode
	}
	defer resp.Body.Close()
	var ipRates ipGeoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&ipRates); err != nil {
		return nil, err
	}
	result := make(map[string]float64)
	for currency, rate := range ipRates.Rates {
		value, err := strconv.ParseFloat(rate.RateForAmount, 64)
		if err != nil {
			return nil, err
		}
		result[currency] = value
	}
	return result, nil
}

type IPCurrencyRates map[string]float64

func (ipcr IPCurrencyRates) Rate(currency string) float64 {
	return ipcr[strings.ToUpper(currency)]
}

type ipGeoAPIResponse struct {
	BaseCurrencyCode string               `json:"base_currency_code"`
	BaseCurrencyName string               `json:"base_currency_name"`
	Amount           string               `json:"amount"`
	UpdatedDate      string               `json:"updated_date"`
	Rates            map[string]ipGeoRate `json:"rates"`
}
type ipGeoRate struct {
	CurrencyName  string `json:"currency_name"`
	Rate          string `json:"rate"`
	RateForAmount string `json:"rate_for_amount"`
}
