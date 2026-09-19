// Package domain defines core business entities and provider interfaces for currencies.
package domain

import (
	"context"
	"errors"
)

var (
	ErrCurrencyNotFound    = errors.New("currency not found in ISO 4217 registry")
	ErrInvalidCurrencyCode = errors.New("invalid ISO 4217 3-letter currency code")
	ErrInvalidAmount       = errors.New("invalid amount for currency formatting")
)

// Currency represents an official ISO 4217 currency with dynamic symbol and decimal precision.
// Memory alignment optimized: 16-byte strings first, followed by 8-byte ints.
type Currency struct {
	Code           string `json:"code"`            // ISO 4217 Alpha Code (e.g. "USD", "EUR", "PEN")
	Symbol         string `json:"symbol"`          // Official currency symbol (e.g. "S/", "$", "€")
	NumericCode    int    `json:"numeric_code"`    // ISO 4217 Numeric Code (e.g. 604, 840, 978)
	FractionDigits int    `json:"fraction_digits"` // Standard minor unit decimal places
}

// CurrencyProvider defines the repository interface for querying official currency data.
type CurrencyProvider interface {
	ListCurrencies(ctx context.Context) ([]Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*Currency, error)
	ValidateCurrencyCode(ctx context.Context, code string) (*Currency, error)
	SearchCurrencies(ctx context.Context, query string) ([]Currency, error)
	FormatAmount(ctx context.Context, code string, amount float64) (string, error)
}
