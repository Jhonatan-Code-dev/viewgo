package domain

import (
	"context"
	"errors"
)

var (
	ErrCountryNotFound = errors.New("country not found")
	ErrInvalidISO2Code = errors.New("invalid ISO 3166-1 alpha-2 country code")
)

// Country represents a sovereign state or territory officially registered in ISO 3166-1 / Unicode CLDR.
type Country struct {
	Alpha2     string `json:"alpha2"`     // e.g. "ES", "US", "CO", "FR", "JP"
	Alpha3     string `json:"alpha3"`     // e.g. "ESP", "USA", "COL", "FRA", "JPN"
	Numeric    int    `json:"numeric"`    // ISO 3166-1 numeric code
	Name       string `json:"name"`       // Localized English / Standard Display Name
	NativeName string `json:"native_name"` // Native Language Display Name if available
	IsOfficial bool   `json:"is_official"` // ISO official state status
}

// CountryProvider defines the repository interface for querying dynamic country datasets.
type CountryProvider interface {
	ListCountries(ctx context.Context) ([]Country, error)
	GetCountryByCode(ctx context.Context, code string) (*Country, error)
	SearchCountries(ctx context.Context, query string) ([]Country, error)
}
