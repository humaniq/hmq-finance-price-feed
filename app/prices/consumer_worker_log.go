package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type LogWorker struct {
}

func (lw *LogWorker) Work(ctx context.Context, values []price.Value) error {
	ctx = context.WithValue(ctx, "tag", "LOG_WORKER")
	for _, value := range values {
		app.Logger().Debug(ctx, "- %+v", value)
	}
	return nil
}
