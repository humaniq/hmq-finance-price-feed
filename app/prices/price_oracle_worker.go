package prices

import (
	"context"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type PriceOracleWriteWorker struct {
}

func NewPriceOracleWriteWorker() *PriceOracleWriteWorker {
	return &PriceOracleWriteWorker{}
}
func (poww *PriceOracleWriteWorker) Work(ctx context.Context, values []price.Value) {

}
