package application

import (
	"context"

	"viewgo/internal/modules/currency/domain"
)

type CurrencyUseCase struct {
	provider domain.CurrencyProvider
}

func NewCurrencyUseCase(provider domain.CurrencyProvider) *CurrencyUseCase {
	return &CurrencyUseCase{
		provider: provider,
	}
}

func (uc *CurrencyUseCase) ListAll(ctx context.Context) ([]domain.Currency, error) {
	return uc.provider.ListCurrencies(ctx)
}

func (uc *CurrencyUseCase) GetByCode(ctx context.Context, code string) (*domain.Currency, error) {
	return uc.provider.GetCurrencyByCode(ctx, code)
}

func (uc *CurrencyUseCase) Search(ctx context.Context, query string) ([]domain.Currency, error) {
	return uc.provider.SearchCurrencies(ctx, query)
}

func (uc *CurrencyUseCase) FormatAmount(ctx context.Context, code string, amount float64) (string, error) {
	return uc.provider.FormatAmount(ctx, code, amount)
}
