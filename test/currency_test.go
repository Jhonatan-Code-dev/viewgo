// Package test provides unified, isolated unit and integration tests for ViewGo.
package test

import (
	"context"
	"testing"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/currency"
)

func TestCurrencyProvider_ListAndLookup(t *testing.T) {
	provider, err := currency.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Currency Provider: %v", err)
	}

	ctx := context.Background()
	currencies, err := provider.ListCurrencies(ctx)
	if err != nil {
		t.Fatalf("ListCurrencies returned error: %v", err)
	}

	if len(currencies) == 0 {
		t.Errorf("Expected non-empty currencies list from ISO 4217 CLDR registry")
	}

	t.Logf("Successfully dynamically loaded %d official ISO 4217 currencies", len(currencies))

	// Dynamically test GetCurrencyByCode and FormatAmount using loaded CLDR currency records
	sampleCount := 10
	if len(currencies) < sampleCount {
		sampleCount = len(currencies)
	}

	for i := 0; i < sampleCount; i++ {
		target := currencies[i]
		c, err := provider.GetCurrencyByCode(ctx, target.Code)
		if err != nil {
			t.Errorf("Failed to find currency by code %s: %v", target.Code, err)
			continue
		}
		if c.Symbol == "" || c.Code != target.Code {
			t.Errorf("Incomplete or mismatched currency data for %s: %+v", target.Code, c)
		}

		formatted, err := provider.FormatAmount(ctx, target.Code, 1234.56)
		if err != nil {
			t.Errorf("Failed to format amount for %s: %v", target.Code, err)
		} else {
			t.Logf("Formatted 1234.56 in %s -> %s", target.Code, formatted)
		}
	}
}

func TestCurrencyProvider_ValidateCurrencyCode(t *testing.T) {
	provider, err := currency.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Currency Provider: %v", err)
	}

	ctx := context.Background()

	// Test valid 3-letter currency codes (e.g. PEN for Peru Sol, USD, EUR, GBP, JPY, CAD)
	validCodes := []string{"PEN", "pen", "USD", "usd", "EUR", "eur", "GBP", "JPY"}
	for _, code := range validCodes {
		c, err := provider.ValidateCurrencyCode(ctx, code)
		if err != nil {
			t.Errorf("Expected valid ISO 4217 currency code for %s, got error: %v", code, err)
		} else if c.Code == "" || c.Symbol == "" {
			t.Errorf("Expected valid Currency struct details for %s, got %+v", code, c)
		}
	}

	// Test Peru specific currency (PEN)
	pen, err := provider.ValidateCurrencyCode(ctx, "PEN")
	if err != nil {
		t.Fatalf("Failed to validate PEN (Peruvian Sol): %v", err)
	}
	if pen.Code != "PEN" {
		t.Errorf("Expected Code PEN, got %s", pen.Code)
	}
	t.Logf("PEN Details: Code=%s, NumericCode=%d, Symbol=%s, FractionDigits=%d", pen.Code, pen.NumericCode, pen.Symbol, pen.FractionDigits)

	// Test invalid code length (2-letter code should be rejected by ValidateCurrencyCode)
	_, err = provider.ValidateCurrencyCode(ctx, "PE")
	if err == nil {
		t.Errorf("Expected error for 2-letter code 'PE' in ValidateCurrencyCode, but got success")
	}

	// Test invalid currency code
	_, err = provider.ValidateCurrencyCode(ctx, "XYZ")
	if err == nil {
		t.Errorf("Expected error for non-existent code 'XYZ' in ValidateCurrencyCode, but got success")
	}
}
