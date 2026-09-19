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
