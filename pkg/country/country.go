// Package country provides dynamic, official ISO 3166-1 country datasets derived from Unicode CLDR.
package country

import (
	"context"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/domain"
	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/infrastructure"
)

// Country represents a sovereign state or territory officially registered in ISO 3166-1.
type Country = domain.Country

// Provider defines the interface for querying dynamic country datasets.
type Provider interface {
	ListCountries(ctx context.Context) ([]Country, error)
	GetCountryByCode(ctx context.Context, code string) (*Country, error)
	SearchCountries(ctx context.Context, query string) ([]Country, error)
}

// NewProvider initializes and returns a new official Unicode CLDR Country Provider.
func NewProvider() (Provider, error) {
	return infrastructure.NewCLDRCountryProvider()
}
