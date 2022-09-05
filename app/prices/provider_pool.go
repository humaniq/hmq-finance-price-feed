package prices

import (
	"context"
	"github.com/humaniq/hmq-finance-price-feed/app"
	"github.com/humaniq/hmq-finance-price-feed/app/price"
)

type ProviderPool struct {
	providers map[string]*Provider
	out       chan []price.Value
}

func NewProviderPool() *ProviderPool {
	return &ProviderPool{
		providers: make(map[string]*Provider),
		out:       make(chan []price.Value),
	}
}
func (pp *ProviderPool) AddProvider(provider *Provider) {
	pp.providers[provider.name] = provider
}
func (pp *ProviderPool) Feed() <-chan []price.Value {
	return pp.out
}

func (pp *ProviderPool) Start(ctx context.Context) error {
	for _, provider := range pp.providers {
		provider.Provide(ctx, pp.out)
	}
	return nil
}
func (pp *ProviderPool) Stop(ctx context.Context) {
	for name, provider := range pp.providers {
		app.Logger().Info(ctx, "Stopping provider %s", name)
		provider.Stop()
	}
}
func (pp *ProviderPool) WaitForDone(ctx context.Context) {
	for name, provider := range pp.providers {
		provider.WaitForDone()
		app.Logger().Info(ctx, "Provider %s stopped", name)
	}
	close(pp.out)
}
