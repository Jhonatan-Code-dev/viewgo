package infrastructure

import (
	"context"
	"testing"
)

func TestCLDRCountryProvider_ListCountries(t *testing.T) {
	provider, err := NewCLDRCountryProvider()
	if err != nil {
		t.Fatalf("Failed to initialize CLDRCountryProvider: %v", err)
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

	// Test specific ISO codes lookup
	testCodes := []string{"US", "ES", "CO", "FR", "JP", "ESP", "USA"}
	for _, code := range testCodes {
		c, err := provider.GetCountryByCode(ctx, code)
		if err != nil {
			t.Errorf("Failed to find official country by code %s: %v", code, err)
			continue
		}
		if c.Name == "" || c.Alpha2 == "" {
			t.Errorf("Incomplete country data for %s: %+v", code, c)
		}
	}
}
