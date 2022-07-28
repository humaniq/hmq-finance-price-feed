package getgeoapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ConversionRateRecord struct {
	CurrencyName  string
	Rate          float64
	RateForAmount float64
}

type ConversionRates struct {
	CurrencyCode string
	CurrencyName string
	UpdateDate   time.Time
	Amount       float64
	Rates        map[string]ConversionRateRecord
}

func GetConversionRates(ctx context.Context, client *http.Client, apiKey string, baseCurrency string, amount float64, toCurrencies ...string) (*ConversionRates, error) {
	url := fmt.Sprintf("https://api.getgeoapi.com/v2/currency/convert?api_key=%s&from=%s&amount=%f&format=json", apiKey, baseCurrency, amount)
	if len(toCurrencies) > 0 {
		url = fmt.Sprintf("%s&toCurrencies=%s", url, strings.Join(toCurrencies, ","))
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrWrongStatus, resp.StatusCode)
	}
	defer resp.Body.Close()
	var body conversionRatesResponseBody
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	updateDate, err := time.Parse("2006-04-02", body.UpdatedDate)
	if err != nil {
		return nil, err
	}
	resultAmount, err := strconv.ParseFloat(body.Amount, 64)
	if err != nil {
		return nil, err
	}
	result := ConversionRates{
		CurrencyCode: body.CurrencyCode,
		CurrencyName: body.CurrencyName,
		UpdateDate:   updateDate,
		Amount:       resultAmount,
		Rates:        make(map[string]ConversionRateRecord),
	}
	for currency, rate := range body.Rates {
		rateValue, err := strconv.ParseFloat(rate.Rate, 64)
		if err != nil {
			return nil, err
		}
		rateForAmountValue, err := strconv.ParseFloat(rate.RateForAmount, 64)
		if err != nil {
			return nil, err
		}
		result.Rates[currency] = ConversionRateRecord{
			CurrencyName:  rate.CurrencyName,
			Rate:          rateValue,
			RateForAmount: rateForAmountValue,
		}
	}
	return &result, nil
}

type conversionRatesResponseBody struct {
	CurrencyCode string `json:"base_currency_code"`
	CurrencyName string `json:"base_currency_name"`
	Amount       string `json:"amount"`
	UpdatedDate  string `json:"updated_date"`
	Rates        map[string]struct {
		CurrencyName  string `json:"currency_name"`
		Rate          string `json:"rate"`
		RateForAmount string `json:"rate_for_amount"`
	} `json:"rates"`
}
