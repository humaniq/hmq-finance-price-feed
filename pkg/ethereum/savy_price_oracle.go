package ethereum

import (
	"context"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/humaniq/hmq-finance-price-feed/pkg/ethereum/contracts"
)

type PriceOracleWriter struct {
	conn       *TransactConnection
	transactor *contracts.ContractsTransactor
}

func NewPriceOracleWriter(conn *TransactConnection, contractAddressHex string) (*PriceOracleWriter, error) {
	ct, err := contracts.NewContractsTransactor(common.HexToAddress(contractAddressHex), conn.ethClient)
	if err != nil {
		return nil, err
	}
	return &PriceOracleWriter{
		transactor: ct,
		conn:       conn,
	}, nil
}
func (pow *PriceOracleWriter) SetDirectPrice(ctx context.Context, addressHex string, value float64, decimals float64) (string, error) {
	opts, err := pow.conn.TxBegin(ctx)
	if err != nil {
		return "", err
	}
	defer pow.conn.TxEnd()
	tx, err := pow.transactor.SetDirectPrice(opts, common.HexToAddress(addressHex), toBigIntPrice(value, decimals))
	if err != nil {
		return "", err
	}
	pow.conn.TxWait(ctx, tx)
	return tx.Hash().String(), nil
}

func toBigIntPrice(value float64, decimals float64) *big.Int {
	return big.NewInt(
		int64(value * math.Pow(10, decimals)),
	)
}
