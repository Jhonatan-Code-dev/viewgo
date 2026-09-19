package application

import (
	"context"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/domain"
)

type CountryUseCase struct {
	provider domain.CountryProvider
}

func NewCountryUseCase(provider domain.CountryProvider) *CountryUseCase {
	return &CountryUseCase{
		provider: provider,
	}
}

func (uc *CountryUseCase) ListAll(ctx context.Context) ([]domain.Country, error) {
	return uc.provider.ListCountries(ctx)
}

func (uc *CountryUseCase) GetByCode(ctx context.Context, code string) (*domain.Country, error) {
	return uc.provider.GetCountryByCode(ctx, code)
}

func (uc *CountryUseCase) Search(ctx context.Context, query string) ([]domain.Country, error) {
	return uc.provider.SearchCountries(ctx, query)
}
