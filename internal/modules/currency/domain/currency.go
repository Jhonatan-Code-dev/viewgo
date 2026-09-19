package domain

import (
	"context"
	"errors"
)

var (
	ErrCurrencyNotFound = errors.New("currency not found in ISO 4217 registry")
	ErrInvalidAmount    = errors.New("invalid amount for currency formatting")
)

// Currency represents an official ISO 4217 currency with dynamic symbols and formatting options.
type Currency struct {
	Code           string `json:"code"`            // ISO 4217 Alpha Code (e.g. "USD", "EUR", "JPY", "GBP", "COP")
	NumericCode    int    `json:"numeric_code"`    // ISO 4217 Numeric Code (e.g. 840, 978, 392)
	Symbol         string `json:"symbol"`          // Official currency symbol (e.g. "$", "€", "¥", "£", "COP$")
	NarrowSymbol   string `json:"narrow_symbol"`   // Narrow currency symbol (e.g. "$", "€", "¥")
	Name           string `json:"name"`            // English Display Name from Unicode CLDR (e.g. "US Dollar")
	FractionDigits int    `json:"fraction_digits"` // Standard minor unit decimal places
}

// CurrencyProvider defines the repository interface for querying official currency data.
type CurrencyProvider interface {
	ListCurrencies(ctx context.Context) ([]Currency, error)
	GetCurrencyByCode(ctx context.Context, code string) (*Currency, error)
	SearchCurrencies(ctx context.Context, query string) ([]Currency, error)
	FormatAmount(ctx context.Context, code string, amount float64) (string, error)
}
