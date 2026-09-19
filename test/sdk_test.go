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

func TestRootSDK_TenantConfigValidation(t *testing.T) {
	ctx := context.Background()

	countryProvider, _ := viewgo.NewCountryProvider()
	tzProvider, _ := viewgo.NewTimezoneProvider()
	currProvider, _ := viewgo.NewCurrencyProvider()

	// Test valid Peru SaaS Tenant configuration: Country="PE", Timezone="America/Lima", Currency="PEN"
	tenantCfg, err := viewgo.ValidateTenantConfig(
		ctx, countryProvider, tzProvider, currProvider,
		"PE", "America/Lima", "PEN",
	)
	if err != nil {
		t.Fatalf("Failed to validate Peru SaaS Tenant Config: %v", err)
	}

	if tenantCfg.CountryCode != "PE" || tenantCfg.Timezone != "America/Lima" || tenantCfg.CurrencyCode != "PEN" {
		t.Errorf("Mismatch in validated TenantConfig codes: %+v", tenantCfg)
	}

	if tenantCfg.Currency.Symbol != "S/" || tenantCfg.Currency.FractionDigits != 2 {
		t.Errorf("Mismatch in Peru currency details: %+v", tenantCfg.Currency)
	}

	t.Logf("Validated SaaS Tenant Config: DB Values: CountryCode=%s, Timezone=%s, CurrencyCode=%s | Render Values: Symbol=%s, Decimals=%d",
		tenantCfg.CountryCode, tenantCfg.Timezone, tenantCfg.CurrencyCode,
		tenantCfg.Currency.Symbol, tenantCfg.Currency.FractionDigits)

	// Test invalid country code in payload
	_, err = viewgo.ValidateTenantConfig(ctx, countryProvider, tzProvider, currProvider, "XX", "America/Lima", "PEN")
	if err == nil {
		t.Errorf("Expected error for invalid country code 'XX'")
	}

	// Test invalid timezone in payload
	_, err = viewgo.ValidateTenantConfig(ctx, countryProvider, tzProvider, currProvider, "PE", "Invalid/Timezone", "PEN")
	if err == nil {
		t.Errorf("Expected error for invalid timezone 'Invalid/Timezone'")
	}

	// Test invalid currency code in payload
	_, err = viewgo.ValidateTenantConfig(ctx, countryProvider, tzProvider, currProvider, "PE", "America/Lima", "INVALID")
	if err == nil {
		t.Errorf("Expected error for invalid currency code 'INVALID'")
	}
}
