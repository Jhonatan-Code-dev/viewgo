# ViewGo

[![Go Reference](https://pkg.go.dev/badge/github.com/Jhonatan-Code-dev/viewgo.svg)](https://pkg.go.dev/github.com/Jhonatan-Code-dev/viewgo)
[![Go CI Pipeline](https://github.com/Jhonatan-Code-dev/viewgo/actions/workflows/ci.yml/badge.svg)](https://github.com/Jhonatan-Code-dev/viewgo/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B%20%7C%201.26%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2B%20Modular%20Monolith-blueviolet)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

ViewGo is an enterprise-grade Go SDK and microservice published on [pkg.go.dev](https://pkg.go.dev/github.com/Jhonatan-Code-dev/viewgo). 

Architected with Clean Architecture principles embedded within a Modular Monolith layout, ViewGo provides dynamic, official reference data for Countries, IANA Timezones, and ISO 4217 Currencies & Symbols without relying on hardcoded static datasets.

---

## Installation & Library Usage

Install ViewGo as a Go module:

```bash
go get github.com/Jhonatan-Code-dev/viewgo@latest
```

### Country Package (`pkg/country`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/country"
)

func main() {
	provider, err := country.NewProvider()
	if err != nil {
		log.Fatal(err)
	}

	c, err := provider.GetCountryByCode(context.Background(), "CO")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s (%s / %s) - M.49: %d\n", c.Name, c.Alpha2, c.Alpha3, c.Numeric)
	// Output: Colombia (CO / COL) - M.49: 170
}
```

### Timezone Package (`pkg/timezone`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"
)

func main() {
	provider, err := timezone.NewProvider()
	if err != nil {
		log.Fatal(err)
	}

	tz, err := provider.GetTimezoneByName(context.Background(), "America/Bogota")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s - Offset: %s (Current Time: %s)\n", tz.IANA, tz.UTCOffset, tz.CurrentTime)
	// Output: America/Bogota - Offset: -05:00
}
```

### Currency Package (`pkg/currency`)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/currency"
)

func main() {
	provider, err := currency.NewProvider()
	if err != nil {
		log.Fatal(err)
	}

	formatted, err := provider.FormatAmount(context.Background(), "EUR", 250.50)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Formatted Amount:", formatted)
	// Output: Formatted Amount: € 250.50
}
```

---

## Core Capabilities

- **Zero Static Datasets**: All reference data is dynamically resolved at runtime from standard Go libraries (`time/tzdata`) and official Unicode CLDR registries (`golang.org/x/text`).
- **Modular Monolith & Clean Architecture**: Domain, Application, Infrastructure, and Delivery layers isolated within independent feature modules (`country`, `timezone`, `currency`).
- **ISO 3166-1 Country Registry**: Dynamic lookup of 285+ countries and territories including ISO Alpha-2, Alpha-3, UN M.49 numeric codes, English and localized native names.
- **IANA Timezone Registry**: Resolution of 123+ canonical IANA timezone identifiers with dynamic UTC offset calculation, DST status detection, and real-time localized timestamps.
- **ISO 4217 Currency Engine**: Resolution of 155+ legal tender currencies with official CLDR symbols, narrow symbols, and decimal precision formatting.
- **HTTP REST Microservice**: Out-of-the-box REST API controller with structured JSON endpoints.

---

## Architecture Diagram

```mermaid
graph TD
    subgraph Public SDK Export Packages
        SDK_C["github.com/Jhonatan-Code-dev/viewgo/pkg/country"]
        SDK_T["github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"]
        SDK_R["github.com/Jhonatan-Code-dev/viewgo/pkg/currency"]
    end

    subgraph Delivery Layer HTTP REST
        HTTP_C[Country Handler]
        HTTP_T[Timezone Handler]
        HTTP_R[Currency Handler]
    end

    subgraph Application Layer UseCases
        UC_C[Country UseCase]
        UC_T[Timezone UseCase]
        UC_R[Currency UseCase]
    end

    subgraph Domain Layer Entities & Interfaces
        D_C[Country Domain]
        D_T[Timezone Domain]
        D_R[Currency Domain]
    end

    subgraph Infrastructure Layer Official Providers
        INF_C[CLDR Country Provider<br/>golang.org/x/text/language]
        INF_T[IANA Timezone Provider<br/>time/tzdata]
        INF_R[CLDR Currency Provider<br/>golang.org/x/text/currency]
    end

    SDK_C --> INF_C
    SDK_T --> INF_T
    SDK_R --> INF_R

    HTTP_C --> UC_C
    HTTP_T --> UC_T
    HTTP_R --> UC_R

    UC_C --> D_C
    UC_T --> D_T
    UC_R --> D_R

    INF_C -. Implements .-> D_C
    INF_T -. Implements .-> D_T
    INF_R -. Implements .-> D_R
```

---

## Directory Layout

```text
viewgo/
├── .agents/                          # AGY Workspace Customization & Architecture Rules
│   ├── AGENTS.md
│   ├── rules/
│   │   ├── golang-clean-architecture.md
│   │   └── no-emojis-strict-documentation.md
│   └── skills/
│       ├── golang-clean-architecture/
│       └── modular-monolith-go/
├── .github/
│   └── workflows/
│       └── ci.yml                    # Automated GitHub Actions & pkg.go.dev Indexing Pipeline
├── cmd/
│   └── viewgo/
│       └── main.go                   # Main REST Server & Dynamic Statistics Executable
├── pkg/                              # Exported Public Go SDK Packages (pkg.go.dev)
│   ├── country/                      # Public Country Package (country.NewProvider)
│   ├── timezone/                     # Public Timezone Package (timezone.NewProvider)
│   └── currency/                     # Public Currency Package (currency.NewProvider)
├── internal/
│   ├── shared/
│   │   └── domain/errors.go          # Shared Kernel Primitives & Sentinel Errors
│   └── modules/
│       ├── country/                  # Country Feature Module (Clean Architecture)
│       ├── timezone/                 # Timezone Feature Module (Clean Architecture)
│       └── currency/                 # Currency Feature Module (Clean Architecture)
├── doc.go                            # Top-level GoDoc package documentation
├── go.mod                            # Go Module definition
├── LICENSE                           # MIT License
└── README.md
```

---

## REST API Specification

### 1. Country Endpoints

| Method | Path | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/countries` | List all countries (Supports query parameter `?q=search`) |
| `GET` | `/api/v1/countries/{code}` | Get country by ISO Alpha-2 (`ES`, `CO`) or Alpha-3 (`COL`, `ESP`) |

#### Example Response (`GET /api/v1/countries/CO`):
```json
{
  "status": "success",
  "data": {
    "alpha2": "CO",
    "alpha3": "COL",
    "numeric": 170,
    "name": "Colombia",
    "native_name": "Colombia",
    "is_official": true
  }
}
```

---

### 2. Timezone Endpoints

| Method | Path | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/timezones` | List all canonical IANA timezones (Supports query parameter `?q=search`) |
| `GET` | `/api/v1/timezones/{ianaName}` | Get timezone details by IANA identifier (`America/Bogota`) |

#### Example Response (`GET /api/v1/timezones/America/Bogota`):
```json
{
  "status": "success",
  "data": {
    "iana": "America/Bogota",
    "abbreviation": "-05",
    "utc_offset": "-05:00",
    "raw_offset_seconds": -18000,
    "is_dst": false,
    "current_time": "2026-09-19T13:46:59-05:00"
  }
}
```

---

### 3. Currency Endpoints

| Method | Path | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/currencies` | List all ISO 4217 legal tender currencies and symbols |
| `GET` | `/api/v1/currencies/{code}` | Get currency by ISO code (`USD`, `EUR`, `JPY`, `COP`) |
| `GET` | `/api/v1/currencies/format?code=EUR&amount=250.50` | Format monetary amount with official currency symbol |

#### Example Response (`GET /api/v1/currencies/format?code=EUR&amount=250.50`):
```json
{
  "status": "success",
  "currency": "EUR",
  "amount": 250.5,
  "formatted": "€ 250.50"
}
```

---

## Release & Versioning Guidelines

To publish a release to `pkg.go.dev`:

1. Tag a semantic version in Git:
   ```bash
   git tag v0.1.0
   git push origin v0.1.0
   ```

2. Trigger Go Proxy Indexing:
   ```bash
   curl https://proxy.golang.org/github.com/Jhonatan-Code-dev/viewgo/@v/v0.1.0.info
   ```

---

## License

This project is licensed under the [MIT License](LICENSE).
Copyright (c) 2026 Jhonatan-Code-dev.
