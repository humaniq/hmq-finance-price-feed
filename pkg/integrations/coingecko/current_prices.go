package coingecko

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type PriceRecord struct {
	Id       string
	Currency string
	Value    float64
}

func GetCurrentPrices(ctx context.Context, client *http.Client, ids []string, currencies []string) ([]PriceRecord, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=%s",
		strings.Join(ids, ","),
		strings.Join(currencies, ","))
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrWrongStatus, resp.StatusCode)
	}
	defer resp.Body.Close()
	body := make(map[string]map[string]float64)
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, err
	}
	list := make([]PriceRecord, 0, len(ids)*len(currencies))
	for id, value := range body {
		for currency, price := range value {
			list = append(list, PriceRecord{
				Id:       id,
				Currency: currency,
				Value:    price,
			})
		}
	}
	return list, nil
}
