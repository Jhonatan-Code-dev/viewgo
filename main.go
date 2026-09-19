package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	countryApp "viewgo/internal/modules/country/application"
	countryDel "viewgo/internal/modules/country/delivery"
	countryInf "viewgo/internal/modules/country/infrastructure"

	tzApp "viewgo/internal/modules/timezone/application"
	tzDel "viewgo/internal/modules/timezone/delivery"
	tzInf "viewgo/internal/modules/timezone/infrastructure"

	currApp "viewgo/internal/modules/currency/application"
	currDel "viewgo/internal/modules/currency/delivery"
	currInf "viewgo/internal/modules/currency/infrastructure"
)

func main() {
	fmt.Println("================================================================================")
	fmt.Println("   WORLD DATA SERVICES - CLEAN ARCHITECTURE & MODULAR MONOLITH (GOLANG)")
	fmt.Println("================================================================================")

	ctx := context.Background()

	// 1. Initialize Country Module
	countryProvider, err := countryInf.NewCLDRCountryProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Country Provider: %v", err)
	}
	countryUseCase := countryApp.NewCountryUseCase(countryProvider)
	countryHandler := countryDel.NewCountryHTTPHandler(countryUseCase)

	// 2. Initialize Timezone Module
	tzProvider, err := tzInf.NewIANATimezoneProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Timezone Provider: %v", err)
	}
	tzUseCase := tzApp.NewTimezoneUseCase(tzProvider)
	tzHandler := tzDel.NewTimezoneHTTPHandler(tzUseCase)

	// 3. Initialize Currency Module
	currProvider, err := currInf.NewCLDRCurrencyProvider()
	if err != nil {
		log.Fatalf("Failed to initialize Currency Provider: %v", err)
	}
	currUseCase := currApp.NewCurrencyUseCase(currProvider)
	currHandler := currDel.NewCurrencyHTTPHandler(currUseCase)

	// Print dynamic dataset stats
	countries, _ := countryUseCase.ListAll(ctx)
	timezones, _ := tzUseCase.ListAll(ctx)
	currencies, _ := currUseCase.ListAll(ctx)

	fmt.Printf("✔ [Module: Country]   Loaded %d official ISO 3166-1 countries/territories dynamically.\n", len(countries))
	fmt.Printf("✔ [Module: Timezone]  Loaded %d official IANA timezones dynamically.\n", len(timezones))
	fmt.Printf("✔ [Module: Currency]  Loaded %d official ISO 4217 currencies and symbols dynamically.\n", len(currencies))
	fmt.Println("--------------------------------------------------------------------------------")

	// Print samples
	if len(countries) > 0 {
		fmt.Printf("  Sample Country:  %s (%s / %s)\n", countries[0].Name, countries[0].Alpha2, countries[0].Alpha3)
	}
	if len(timezones) > 0 {
		fmt.Printf("  Sample Timezone: %s (Offset: %s, Time: %s)\n", timezones[0].IANA, timezones[0].UTCOffset, timezones[0].CurrentTime)
	}
	if len(currencies) > 0 {
		sampleFmt, _ := currUseCase.FormatAmount(ctx, "USD", 1500.75)
		fmt.Printf("  Sample Currency: %s (%s) -> Format 1500.75: %s\n", currencies[0].Name, currencies[0].Symbol, sampleFmt)
	}

	fmt.Println("================================================================================")

	// Register REST API Routes
	mux := http.NewServeMux()
	countryHandler.RegisterRoutes(mux)
	tzHandler.RegisterRoutes(mux)
	currHandler.RegisterRoutes(mux)

	// Root status endpoint
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"service":"viewgo-world-data","status":"active","modules":["country","timezone","currency"]}`)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	fmt.Printf("🚀 REST API Server running on http://localhost:%s\n", port)
	fmt.Println("   Endpoints:")
	fmt.Println("     - GET /api/v1/countries")
	fmt.Println("     - GET /api/v1/countries/ES (or /US, /CO, /FRA)")
	fmt.Println("     - GET /api/v1/timezones")
	fmt.Println("     - GET /api/v1/timezones/America/Bogota")
	fmt.Println("     - GET /api/v1/currencies")
	fmt.Println("     - GET /api/v1/currencies/format?code=EUR&amount=250.50")
	fmt.Println("================================================================================")

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("HTTP server error: %v", err)
	}
}
