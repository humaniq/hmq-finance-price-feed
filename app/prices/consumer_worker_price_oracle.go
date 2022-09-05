package prices

import (
	"context"
	"errors"
	"fmt"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/config"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/ethereum"
)

type PriceOracleWriteWorker struct {
	writer  *ethereum.PriceOracleWriter
	symbols map[string]config.EthNetworkSymbolContract
}

func NewPriceOracleWriteWorker(symbols map[string]config.EthNetworkSymbolContract, writer *ethereum.PriceOracleWriter) *PriceOracleWriteWorker {
	return &PriceOracleWriteWorker{
		writer:  writer,
		symbols: symbols,
	}
}
func (poww *PriceOracleWriteWorker) Work(ctx context.Context, values []price.Value) error {
	for _, value := range values {
		symbol, found := poww.symbols[value.Symbol]
		if !found {
			return errors.New(fmt.Sprintf("symbol address not found for %s", value.Symbol))
		}
		txId, err := poww.writer.SetDirectPrice(ctx, symbol.AddressHex, value.Price, symbol.Decimals)
		if err != nil {
			return err
		}
		app.Logger().Info(context.WithValue(ctx, "TxID", txId), "[PriceOracleWorker] %f price written for %s=[%s]", value.Price, value.Symbol, symbol.AddressHex)
	}
	return nil
}
