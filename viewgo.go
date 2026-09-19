// Package viewgo provides high-performance, dynamic Go SDK providers for
// official Countries (ISO 3166-1), IANA Timezones, and Currencies & Symbols (ISO 4217).
//
// Designed as a pure Go library module for embedding into Go applications without
// external network HTTP overhead or hardcoded static datasets.
package viewgo

import (
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
