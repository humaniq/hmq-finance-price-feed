package prices

import (
	"context"
	"time"

	"github.com/humaniq/hmq-finance-price-feed/app/price"
	"github.com/humaniq/hmq-finance-price-feed/pkg/logger"
)

type ProviderFunc func(ctx context.Context) ([]price.Value, error)

type Provider struct {
	fn     ProviderFunc
	name   string
	done   chan interface{}
	ticker *time.Ticker
	every  time.Duration
}

func NewProvider(name string, fn ProviderFunc, every time.Duration) *Provider {
	return &Provider{
		name:  name,
		fn:    fn,
		every: every,
	}
}

func (p *Provider) Name() string {
	return p.name
}

func (p *Provider) Provide(ctx context.Context, out chan<- []price.Value) error {
	logger.Info(ctx, "Provider %s, every %v", p.Name(), p.every)
	p.ticker = time.NewTicker(p.every)
	go p.Run(ctx, p.ticker, out)
	return nil
}
func (p *Provider) Stop() {
	p.ticker.Stop()
}
func (p *Provider) WaitForDone() {
	<-p.done
}

func (p *Provider) Run(ctx context.Context, ticker *time.Ticker, out chan<- []price.Value) {
	defer close(p.done)
	for range ticker.C {
		values, err := p.fn(ctx)
		if err != nil {
			logger.Error(ctx, "ERROR GETTING COINGECKO DATA FOR %s: %s", p.name, err)
			continue
		}
		out <- values
	}
}
