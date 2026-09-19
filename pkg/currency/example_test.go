package currency_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/currency"
)

func ExampleNewProvider() {
	provider, err := currency.NewProvider()
	if err != nil {
		log.Fatalf("Failed to initialize currency provider: %v", err)
	}

	ctx := context.Background()
	formatted, err := provider.FormatAmount(ctx, "EUR", 250.50)
	if err != nil {
		log.Fatalf("Formatting error: %v", err)
	}

	fmt.Println(formatted)
	// Output:
	// € 250.50
}
