// Package currency provides dynamic, official ISO 4217 currencies and symbols derived from Unicode CLDR.
package currency

import (
	"context"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/domain"
	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/infrastructure"
)

// Currency represents an official ISO 4217 currency.
type Currency = domain.Currency

// Provider defines the interface for querying dynamic ISO 4217 currencies.
type Provider interface {
	ListCurrencies(ctx context.Context) ([]Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*Currency, error)
	SearchCurrencies(ctx context.Context, query string) ([]Currency, error)
	FormatAmount(ctx context.Context, code string, amount float64) (string, error)
}

// NewProvider initializes and returns a new official Unicode CLDR Currency Provider.
func NewProvider() (Provider, error) {
	return infrastructure.NewCLDRCurrencyProvider()
}
