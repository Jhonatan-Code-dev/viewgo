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
