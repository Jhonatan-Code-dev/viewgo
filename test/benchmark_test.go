package test

import (
	"context"
	"testing"

	"github.com/Jhonatan-Code-dev/viewgo"
	"github.com/Jhonatan-Code-dev/viewgo/pkg/country"
	"github.com/Jhonatan-Code-dev/viewgo/pkg/currency"
	"github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"
)

func BenchmarkCountryLookup_Alpha2(b *testing.B) {
	provider, err := country.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Country Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.GetCountryByCode(ctx, "CO")
	}
}

func BenchmarkCountry_ValidateAlpha2(b *testing.B) {
	provider, err := country.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Country Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.ValidateAlpha2(ctx, "PE")
	}
}

func BenchmarkTimezoneLookup_IANA(b *testing.B) {
	provider, err := timezone.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Timezone Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.GetTimezoneByName(ctx, "America/Bogota")
	}
}

func BenchmarkTimezone_ValidateIANA(b *testing.B) {
	provider, err := timezone.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Timezone Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.ValidateIANAZone(ctx, "America/Lima")
	}
}

func BenchmarkCurrencyLookup_Code(b *testing.B) {
	provider, err := currency.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Currency Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.GetCurrencyByCode(ctx, "EUR")
	}
}

func BenchmarkCurrency_ValidateCode(b *testing.B) {
	provider, err := currency.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Currency Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.ValidateCurrencyCode(ctx, "PEN")
	}
}

func BenchmarkCurrencyFormat_Amount(b *testing.B) {
	provider, err := currency.NewProvider()
	if err != nil {
		b.Fatalf("Failed to initialize Currency Provider: %v", err)
	}

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = provider.FormatAmount(ctx, "EUR", 250.50)
	}
}

func BenchmarkRootSDK_AllLookups(b *testing.B) {
	countryP, _ := viewgo.NewCountryProvider()
	tzP, _ := viewgo.NewTimezoneProvider()
	currP, _ := viewgo.NewCurrencyProvider()

	ctx := context.Background()
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = countryP.GetCountryByCode(ctx, "ES")
		_, _ = tzP.GetTimezoneByName(ctx, "Europe/London")
		_, _ = currP.FormatAmount(ctx, "USD", 1500.75)
	}
}
