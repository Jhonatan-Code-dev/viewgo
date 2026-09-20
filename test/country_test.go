// Package test provides unified, isolated unit and integration tests for ViewGo.
package test

import (
	"context"
	"testing"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/country"
)

func TestCountryProvider_ListAndLookup(t *testing.T) {
	provider, err := country.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Country Provider: %v", err)
	}

	ctx := context.Background()
	countries, err := provider.ListCountries(ctx)
	if err != nil {
		t.Fatalf("ListCountries returned error: %v", err)
	}

	if len(countries) == 0 {
		t.Errorf("Expected non-empty countries list from Unicode CLDR registry")
	}

	t.Logf("Successfully dynamically loaded %d official countries/territories", len(countries))

	// Dynamically test GetCountryByCode using loaded CLDR records
	sampleCount := 10
	if len(countries) < sampleCount {
		sampleCount = len(countries)
	}

	for i := 0; i < sampleCount; i++ {
		target := countries[i]

		byAlpha2, err := provider.GetCountryByCode(ctx, target.Alpha2)
		if err != nil {
			t.Errorf("Failed to find country by Alpha2 code %s: %v", target.Alpha2, err)
			continue
		}
		if byAlpha2.Name == "" || byAlpha2.Alpha2 != target.Alpha2 {
			t.Errorf("Mismatch or incomplete data for Alpha2 %s: %+v", target.Alpha2, byAlpha2)
		}

		if target.Alpha3 != "" {
			byAlpha3, err := provider.GetCountryByCode(ctx, target.Alpha3)
			if err != nil {
				t.Errorf("Failed to find country by Alpha3 code %s: %v", target.Alpha3, err)
				continue
			}
			if byAlpha3.Alpha2 != target.Alpha2 {
				t.Errorf("Mismatch between Alpha2 and Alpha3 lookup for %s", target.Alpha3)
			}
		}
	}
}

func TestCountryProvider_ValidateAlpha2(t *testing.T) {
	provider, err := country.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Country Provider: %v", err)
	}

	ctx := context.Background()

	// Test valid 2-letter codes (e.g., PE for Peru, CO for Colombia, ES for Spain, US for United States)
	validCodes := []string{"PE", "pe", "CO", "co", "ES", "US"}
	for _, code := range validCodes {
		c, err := provider.ValidateAlpha2(ctx, code)
		if err != nil {
			t.Errorf("Expected valid Alpha-2 code for %s, got error: %v", code, err)
		} else if c.Alpha2 == "" {
			t.Errorf("Expected non-empty country struct for %s", code)
		}
	}

	// Test Peru specific verification
	peru, err := provider.ValidateAlpha2(ctx, "PE")
	if err != nil {
		t.Fatalf("Failed to validate PE (Peru): %v", err)
	}
	if peru.Alpha2 != "PE" {
		t.Errorf("Expected Alpha2 PE, got %s", peru.Alpha2)
	}

	// Test invalid code length (3-letter alpha-3 code should be rejected by ValidateAlpha2)
	_, err = provider.ValidateAlpha2(ctx, "PER")
	if err == nil {
		t.Errorf("Expected error for 3-letter code 'PER' in ValidateAlpha2, but got success")
	}

	// Test invalid 2-letter code
	_, err = provider.ValidateAlpha2(ctx, "XX")
	if err == nil {
		t.Errorf("Expected error for non-existent code 'XX' in ValidateAlpha2, but got success")
	}
}

func TestCountryProvider_ValidateAlpha3(t *testing.T) {
	provider, err := country.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Country Provider: %v", err)
	}

	ctx := context.Background()

	// Test valid 3-letter codes (e.g., PER for Peru, COL for Colombia, ESP for Spain, USA for United States)
	validCodes := []string{"PER", "per", "COL", "col", "ESP", "USA"}
	for _, code := range validCodes {
		c, err := provider.ValidateAlpha3(ctx, code)
		if err != nil {
			t.Errorf("Expected valid Alpha-3 code for %s, got error: %v", code, err)
		} else if c.Alpha3 == "" {
			t.Errorf("Expected non-empty country struct for %s", code)
		}
	}

	// Test Peru specific verification (PER)
	peru, err := provider.ValidateAlpha3(ctx, "PER")
	if err != nil {
		t.Fatalf("Failed to validate PER (Peru): %v", err)
	}
	if peru.Alpha3 != "PER" {
		t.Errorf("Expected Alpha3 PER, got %s", peru.Alpha3)
	}

	// Test invalid code length (2-letter alpha-2 code should be rejected by ValidateAlpha3)
	_, err = provider.ValidateAlpha3(ctx, "PE")
	if err == nil {
		t.Errorf("Expected error for 2-letter code 'PE' in ValidateAlpha3, but got success")
	}

	// Test invalid 3-letter code
	_, err = provider.ValidateAlpha3(ctx, "XYZ")
	if err == nil {
		t.Errorf("Expected error for non-existent code 'XYZ' in ValidateAlpha3, but got success")
	}
}
