package infrastructure

import (
	"context"
	"testing"
)

func TestIANATimezoneProvider_ListTimezones(t *testing.T) {
	provider, err := NewIANATimezoneProvider()
	if err != nil {
		t.Fatalf("Failed to initialize IANATimezoneProvider: %v", err)
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

	// Test specific IANA lookups
	testZones := []string{"America/New_York", "Europe/London", "Asia/Tokyo", "America/Bogota", "UTC"}
	for _, zoneName := range testZones {
		tz, err := provider.GetTimezoneByName(ctx, zoneName)
		if err != nil {
			t.Errorf("Failed to find official timezone %s: %v", zoneName, err)
			continue
		}
		if tz.UTCOffset == "" || tz.IANA == "" {
			t.Errorf("Incomplete timezone data for %s: %+v", zoneName, tz)
		}
	}
}
