// Package test provides unified, isolated unit and integration tests for ViewGo.
package test

import (
	"context"
	"testing"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"
)

func TestTimezoneProvider_ListAndLookup(t *testing.T) {
	provider, err := timezone.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Timezone Provider: %v", err)
	}

	ctx := context.Background()
	zones, err := provider.ListTimezones(ctx)
	if err != nil {
		t.Fatalf("ListTimezones returned error: %v", err)
	}

	if len(zones) == 0 {
		t.Errorf("Expected non-empty timezones list from standard tzdata database")
	}

	t.Logf("Successfully dynamically loaded %d official IANA timezones", len(zones))

	// Dynamically test GetTimezoneByName using loaded IANA records
	sampleCount := 10
	if len(zones) < sampleCount {
		sampleCount = len(zones)
	}

	for i := 0; i < sampleCount; i++ {
		target := zones[i]
		tz, err := provider.GetTimezoneByName(ctx, target.IANA)
		if err != nil {
			t.Errorf("Failed to find timezone by IANA name %s: %v", target.IANA, err)
			continue
		}
		if tz.UTCOffset == "" || tz.IANA != target.IANA {
			t.Errorf("Incomplete or mismatched timezone data for %s: %+v", target.IANA, tz)
		}
	}
}

func TestTimezoneProvider_ValidateIANAZone(t *testing.T) {
	provider, err := timezone.NewProvider()
	if err != nil {
		t.Fatalf("Failed to initialize Timezone Provider: %v", err)
	}

	ctx := context.Background()

	// Test valid IANA zones (e.g., America/Lima for Peru, America/Bogota, Europe/London, Asia/Tokyo, UTC)
	validZones := []string{"America/Lima", "America/Bogota", "Europe/London", "Asia/Tokyo", "UTC", "Etc/UTC"}
	for _, z := range validZones {
		tz, err := provider.ValidateIANAZone(ctx, z)
		if err != nil {
			t.Errorf("Expected valid IANA zone for %s, got error: %v", z, err)
		} else if tz.IANA == "" || tz.UTCOffset == "" {
			t.Errorf("Expected valid Timezone struct details for %s, got %+v", z, tz)
		}
	}

	// Test Peru specific IANA timezone (America/Lima)
	lima, err := provider.ValidateIANAZone(ctx, "America/Lima")
	if err != nil {
		t.Fatalf("Failed to validate America/Lima: %v", err)
	}
	if lima.IANA != "America/Lima" {
		t.Errorf("Expected IANA America/Lima, got %s", lima.IANA)
	}

	// Test invalid / non-canonical zone name
	_, err = provider.ValidateIANAZone(ctx, "Invalid/City_Name")
	if err == nil {
		t.Errorf("Expected error for non-existent IANA zone 'Invalid/City_Name', but got success")
	}

	// Test invalid prefix
	_, err = provider.ValidateIANAZone(ctx, "FakePrefix/Bogota")
	if err == nil {
		t.Errorf("Expected error for invalid prefix zone 'FakePrefix/Bogota', but got success")
	}
}
