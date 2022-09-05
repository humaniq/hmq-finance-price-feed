package pancakeswap

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type TokenData struct {
	TimeStamp time.Time
	Name      string
	Symbol    string
	PriceUsd  float64
	PriceBNB  float64
}

func V2Token(ctx context.Context, addressHex string) (*TokenData, error) {
	request, err := http.NewRequest("GET", fmt.Sprintf("https://api.pancakeswap.info/api/v2/tokens/%s", addressHex), nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequestInvalid, err)
	}
	requestCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()
	request = request.WithContext(requestCtx)

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrRequestFailed, err)
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status bad (%d)", ErrRequestFailed, response.StatusCode)
	}
	defer response.Body.Close()

	var responseBody v2ApiResponse[apiTokenData]
	if err := json.NewDecoder(response.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrResponseInvalid, err)
	}

	priceUsd, err := strconv.ParseFloat(responseBody.Data.PriceUSD, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing usd error %s", ErrResponseInvalid, err)
	}
	priceBnb, err := strconv.ParseFloat(responseBody.Data.PriceBNB, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: parsing bnb error %s", ErrResponseInvalid, err)
	}
	return &TokenData{
		TimeStamp: time.Unix(responseBody.UpdatedAt/1000, responseBody.UpdatedAt%1000),
		Name:      responseBody.Data.Name,
		Symbol:    responseBody.Data.Symbol,
		PriceUsd:  priceUsd,
		PriceBNB:  priceBnb,
	}, nil
}

type apiTokenData struct {
	Name     string `json:"name,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	PriceUSD string `json:"price"`
	PriceBNB string `json:"price_BNB"`
}
