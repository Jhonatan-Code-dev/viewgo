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
