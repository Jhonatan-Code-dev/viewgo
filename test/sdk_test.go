// Package test provides unified, isolated unit and integration tests for ViewGo.
package test

import (
	"context"
	"testing"

	"github.com/Jhonatan-Code-dev/viewgo"
)

func TestRootSDK_InitializationAndQuery(t *testing.T) {
	ctx := context.Background()

	countryProvider, err := viewgo.NewCountryProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Root Country Provider: %v", err)
	}

	tzProvider, err := viewgo.NewTimezoneProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Root Timezone Provider: %v", err)
	}

	currProvider, err := viewgo.NewCurrencyProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Root Currency Provider: %v", err)
	}

	// 1. Verify Country
	colombia, err := countryProvider.GetCountryByCode(ctx, "CO")
	if err != nil || colombia.Alpha2 != "CO" {
		t.Errorf("Root SDK failed country query for CO: %v", err)
	}

	// 2. Verify Timezone
	bogota, err := tzProvider.GetTimezoneByName(ctx, "America/Bogota")
	if err != nil || bogota.IANA != "America/Bogota" {
		t.Errorf("Root SDK failed timezone query for America/Bogota: %v", err)
	}

	// 3. Verify Currency
	formatted, err := currProvider.FormatAmount(ctx, "EUR", 250.50)
	if err != nil || formatted == "" {
		t.Errorf("Root SDK failed currency formatting for EUR: %v", err)
	}
}
