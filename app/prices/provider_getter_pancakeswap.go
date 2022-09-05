package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/pkg/integrations/pancakeswap"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type PancakeSwap struct {
	symbols map[string]string
}

func NewPancakeSwap(cfg *config.PancakeSwapProvider, network *config.EthNetwork) *PancakeSwap {
	symbols := make(map[string]string)
	for _, smb := range cfg.Symbols {
		if val, found := network.Symbols[smb]; found {
			symbols[val.AddressHex] = smb
		}
	}
	return &PancakeSwap{symbols: symbols}
}
func (ps *PancakeSwap) GetterFunc() ProviderFunc {
	return func(ctx context.Context) ([]price.Value, error) {
		var values []price.Value
		for hex, smb := range ps.symbols {
			rates, err := pancakeswap.V2Token(ctx, hex)
			if err != nil {
				return nil, err
			}
			values = append(values, price.Value{
				TimeStamp: time.Now(),
				Source:    "pankaceswap",
				Symbol:    smb,
				Currency:  "usd",
				Price:     rates.PriceUsd,
			})
			values = append(values, price.Value{
				TimeStamp: time.Now(),
				Source:    "pankaceswap",
				Symbol:    smb,
				Currency:  "bnb",
				Price:     rates.PriceBNB,
			})
		}
		return values, nil
	}
}
