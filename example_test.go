package viewgo_test

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo"
)

func ExampleNewCountryProvider() {
	provider, err := viewgo.NewCountryProvider()
	if err != nil {
		log.Fatal(err)
	}

	c, err := provider.GetCountryByCode(context.Background(), "CO")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s / %s)\n", c.Name, c.Alpha2, c.Alpha3)
	// Output:
	// Colombia (CO / COL)
}

func ExampleNewTimezoneProvider() {
	provider, err := viewgo.NewTimezoneProvider()
	if err != nil {
		log.Fatal(err)
	}

	tz, err := provider.GetTimezoneByName(context.Background(), "America/Bogota")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s)\n", tz.IANA, tz.UTCOffset)
	// Output:
	// America/Bogota (-05:00)
}

func ExampleNewCurrencyProvider() {
	provider, err := viewgo.NewCurrencyProvider()
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := provider.FormatAmount(context.Background(), "EUR", 250.50)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(formatted)
	// Output:
	// € 250.50
}
