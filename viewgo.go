// Package viewgo provides high-performance, dynamic Go SDK providers for
// official Countries (ISO 3166-1), IANA Timezones, and Currencies & Symbols (ISO 4217).
//
// Designed as a pure Go library module for embedding into Go applications without
// external network HTTP overhead or hardcoded static datasets.
package viewgo

import (
	"context"

	countryDom "github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/domain"
	countryInf "github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/infrastructure"

	tzDom "github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/domain"
	tzInf "github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/infrastructure"

	currDom "github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/domain"
	currInf "github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/infrastructure"
)

// Type aliases for core domain models.
type Country = countryDom.Country
type Timezone = tzDom.Timezone
type Currency = currDom.Currency

// Interfaces for providers.
type CountryProvider = countryDom.CountryProvider
type TimezoneProvider = tzDom.TimezoneProvider
type CurrencyProvider = currDom.CurrencyProvider

// NewCountryProvider initializes and returns a dynamic ISO 3166-1 Country Provider.
func NewCountryProvider() (CountryProvider, error) {
	return countryInf.NewCLDRCountryProvider()
}

// NewTimezoneProvider initializes and returns a dynamic IANA Timezone Provider.
func NewTimezoneProvider() (TimezoneProvider, error) {
	return tzInf.NewIANATimezoneProvider()
}

// NewCurrencyProvider initializes and returns a dynamic ISO 4217 Currency Provider.
func NewCurrencyProvider() (CurrencyProvider, error) {
	return currInf.NewCLDRCurrencyProvider()
}

// TenantConfig represents a validated SaaS tenant regional configuration payload.
type TenantConfig struct {
	CountryCode  string    `json:"country_code"`  // ISO 3166-1 Alpha-2 (e.g. "PE", "US", "ES")
	Timezone     string    `json:"timezone"`      // IANA Zone (e.g. "America/Lima", "America/New_York")
	CurrencyCode string    `json:"currency_code"` // ISO 4217 Alpha-3 (e.g. "PEN", "USD", "EUR")
	Country      *Country  `json:"country,omitempty"`
	TimezoneData *Timezone `json:"timezone_data,omitempty"`
	Currency     *Currency `json:"currency,omitempty"`
}

// ValidateTenantConfig validates a complete SaaS tenant registration/onboarding payload in one call.
// Validates ISO 3166-1 alpha-2 country, IANA timezone, and ISO 4217 currency.
func ValidateTenantConfig(
	ctx context.Context,
	countryP CountryProvider,
	tzP TimezoneProvider,
	currP CurrencyProvider,
	countryAlpha2, ianaZone, currencyAlpha3 string,
) (*TenantConfig, error) {
	country, err := countryP.ValidateAlpha2(ctx, countryAlpha2)
	if err != nil {
		return nil, err
	}

	tz, err := tzP.ValidateIANAZone(ctx, ianaZone)
	if err != nil {
		return nil, err
	}

	curr, err := currP.ValidateCurrencyCode(ctx, currencyAlpha3)
	if err != nil {
		return nil, err
	}

	return &TenantConfig{
		CountryCode:  country.Alpha2,
		Timezone:     tz.IANA,
		CurrencyCode: curr.Code,
		Country:      country,
		TimezoneData: tz,
		Currency:     curr,
	}, nil
}
