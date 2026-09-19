package application

import (
	"context"
	"time"

	"viewgo/internal/modules/timezone/domain"
)

type TimezoneUseCase struct {
	provider domain.TimezoneProvider
}

func NewTimezoneUseCase(provider domain.TimezoneProvider) *TimezoneUseCase {
	return &TimezoneUseCase{
		provider: provider,
	}
}

func (uc *TimezoneUseCase) ListAll(ctx context.Context) ([]domain.Timezone, error) {
	return uc.provider.ListTimezones(ctx)
}

func (uc *TimezoneUseCase) GetByName(ctx context.Context, name string) (*domain.Timezone, error) {
	return uc.provider.GetTimezoneByName(ctx, name)
}

func (uc *TimezoneUseCase) Search(ctx context.Context, query string) ([]domain.Timezone, error) {
	return uc.provider.SearchTimezones(ctx, query)
}

func (uc *TimezoneUseCase) GetCurrentTime(ctx context.Context, name string) (time.Time, error) {
	return uc.provider.GetTime(ctx, name)
}
