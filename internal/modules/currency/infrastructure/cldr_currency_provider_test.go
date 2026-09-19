package infrastructure

import (
	"context"
	"testing"
)

func TestCLDRCurrencyProvider_ListCurrencies(t *testing.T) {
	provider, err := NewCLDRCurrencyProvider()
	if err != nil {
		t.Fatalf("Failed to initialize CLDRCurrencyProvider: %v", err)
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

	// Test specific currency lookups and formatting
	testCodes := []string{"USD", "EUR", "JPY", "GBP", "COP"}
	for _, code := range testCodes {
		c, err := provider.GetCurrencyByCode(ctx, code)
		if err != nil {
			t.Errorf("Failed to find official currency %s: %v", code, err)
			continue
		}
		if c.Symbol == "" || c.Name == "" {
			t.Errorf("Incomplete currency data for %s: %+v", code, c)
		}

		formatted, err := provider.FormatAmount(ctx, code, 1234.56)
		if err != nil {
			t.Errorf("Failed to format amount for %s: %v", code, err)
		} else {
			t.Logf("Formatted 1234.56 in %s -> %s", code, formatted)
		}
	}
}
