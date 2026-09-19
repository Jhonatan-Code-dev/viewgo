package timezone_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"
)

func ExampleNewProvider() {
	provider, err := timezone.NewProvider()
	if err != nil {
		log.Fatalf("Failed to initialize timezone provider: %v", err)
	}

	ctx := context.Background()
	tz, err := provider.GetTimezoneByName(ctx, "America/Bogota")
	if err != nil {
		log.Fatalf("Timezone not found: %v", err)
	}

	fmt.Printf("%s (Offset: %s)\n", tz.IANA, tz.UTCOffset)
	// Output:
	// America/Bogota (Offset: -05:00)
}
