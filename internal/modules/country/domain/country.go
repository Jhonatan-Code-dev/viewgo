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
// Struct alignment optimized: 16-byte strings first, followed by 8-byte int and 1-byte bool.
type Country struct {
	Name       string `json:"name"`        // Official localized display name
	NativeName string `json:"native_name"` // Native language display name
	Alpha2     string `json:"alpha2"`      // ISO 3166-1 alpha-2 code (e.g. "ES", "US", "CO")
	Alpha3     string `json:"alpha3"`      // ISO 3166-1 alpha-3 code (e.g. "ESP", "USA", "COL")
	Numeric    int    `json:"numeric"`     // ISO 3166-1 numeric M.49 code
	IsOfficial bool   `json:"is_official"` // ISO official state status
}

// CountryProvider defines the repository interface for querying dynamic country datasets.
type CountryProvider interface {
	ListCountries(ctx context.Context) ([]Country, error)
	GetCountryByCode(ctx context.Context, code string) (*Country, error)
	SearchCountries(ctx context.Context, query string) ([]Country, error)
}
