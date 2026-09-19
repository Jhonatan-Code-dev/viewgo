package country_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/country"
)

func ExampleNewProvider() {
	provider, err := country.NewProvider()
	if err != nil {
		log.Fatalf("Failed to initialize country provider: %v", err)
	}

	ctx := context.Background()
	c, err := provider.GetCountryByCode(ctx, "CO")
	if err != nil {
		log.Fatalf("Country not found: %v", err)
	}

	fmt.Printf("%s (%s / %s)\n", c.Name, c.Alpha2, c.Alpha3)
	// Output:
	// Colombia (CO / COL)
}
