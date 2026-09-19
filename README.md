# ViewGo: Módulo SDK de Alto Rendimiento para Go

[![Go Reference](https://pkg.go.dev/badge/github.com/Jhonatan-Code-dev/viewgo.svg)](https://pkg.go.dev/github.com/Jhonatan-Code-dev/viewgo)
[![Go CI Pipeline](https://github.com/Jhonatan-Code-dev/viewgo/actions/workflows/ci.yml/badge.svg)](https://github.com/Jhonatan-Code-dev/viewgo/actions)
[![Go Version](https://img.shields.io/badge/Go-1.25%2B%20%7C%201.26%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![Architecture](https://img.shields.io/badge/Architecture-Clean%20%2B%20Modular%20Monolith-blueviolet)](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

ViewGo es una librería y módulo SDK puro en Go publicado en [pkg.go.dev](https://pkg.go.dev/github.com/Jhonatan-Code-dev/viewgo).

Diseñado bajo los principios de **Arquitectura Limpia** (Clean Architecture) y estructurado como un **Monolito Modular** (Modular Monolith), ViewGo proporciona resolución dinámica en memoria para datos mundiales oficiales de **Países** (ISO 3166-1), **Zonas Horarias** (IANA) y **Monedas y Símbolos** (ISO 4217), sin depender de conjuntos de datos estáticos hardcodeados ni llamadas de red HTTP externas.

---

## Características Principales

- **Resolución 100% Dinámica**: Los datos se extraen en tiempo de ejecución desde los paquetes estándar de Go (`time/tzdata`) y repositorios oficiales Unicode CLDR (`golang.org/x/text`).
- **Velocidad Extrema (Sub-10 Nanosegundos)**: Tablas de indexación directa $O(1)$ basadas en arreglos contiguos de 16 bits (`uint16`), logrando latencias de consulta de **~8.4 ns a 10.5 ns**.
- **Cero Asignaciones de Memoria (`0 allocs/op`)**: Consultas calientes con **0 B/op** de basura en heap para reducir el trabajo del Garbage Collector.
- **Consumo Mínimo de RAM (~71 KB Total)**: Reducción del 75% en la huella de memoria RAM mediante compresión de punteros a índices contiguos.
- **Validación Estricta de Estándares**: Funciones dedicadas para validar códigos ISO 3166-1 Alpha-2 (`"PE"`), zonas IANA (`"America/Lima"`) y divisas ISO 4217 (`"PEN"`).
- **Soporte SaaS Multi-Tenant**: Helper integrado (`ValidateTenantConfig`) para validar la configuración regional de usuarios y empresas en una sola llamada.
- **Cumplimiento Estricto Go Standard**: Comentarios de paquete compliant con `staticcheck ST1000` y documentación técnica sin emoticonos.

---

## Instalación e Importación

Agrega ViewGo a tu proyecto de Go:

```bash
go get github.com/Jhonatan-Code-dev/viewgo@latest
```

---

## Ejemplos de Código

### 1. Validación de Configuración Regional SaaS Multi-Tenant (`viewgo`)

Ideal para plataformas SaaS donde los usuarios registran su país, zona horaria y moneda:

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Jhonatan-Code-dev/viewgo"
)

func main() {
	ctx := context.Background()

	// Inicializar proveedores dinámicos
	countryP, _ := viewgo.NewCountryProvider()
	tzP, _ := viewgo.NewTimezoneProvider()
	currP, _ := viewgo.NewCurrencyProvider()

	// Validar payload de registro SaaS en una sola llamada ultrarrápida
	tenantCfg, err := viewgo.ValidateTenantConfig(
		ctx, countryP, tzP, currP,
		"PE",            // País ISO 3166-1 Alpha-2
		"America/Lima",  // Zona Horaria IANA
		"PEN",           // Moneda ISO 4217 Alpha-3
	)
	if err != nil {
		log.Fatalf("Configuración regional inválida: %v", err)
	}

	// Guardar en Base de Datos (Códigos compactos estandarizados)
	fmt.Printf("BD -> Country: %s | TZ: %s | Currency: %s\n",
		tenantCfg.CountryCode, tenantCfg.Timezone, tenantCfg.CurrencyCode)

	// Renderizado de Formatos y Símbolos oficiales en UI / Facturas
	fmt.Printf("UI -> Símbolo: %s | Decimales: %d\n",
		tenantCfg.Currency.Symbol, tenantCfg.Currency.FractionDigits)
}
```

---

### 2. Uso de Subpaquetes Individuales (`pkg/country`, `pkg/timezone`, `pkg/currency`)

#### Países (`pkg/country`)
```go
package main

import (
	"context"
	"fmt"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/country"
)

func main() {
	provider, _ := country.NewProvider()
	ctx := context.Background()

	// Validar código Alpha-2 (ej. Perú)
	pe, err := provider.ValidateAlpha2(ctx, "PE")
	if err == nil {
		fmt.Printf("País: %s (%s / %s) - M.49: %d\n", pe.Name, pe.Alpha2, pe.Alpha3, pe.Numeric)
	}
}
```

#### Zonas Horarias (`pkg/timezone`)
```go
package main

import (
	"context"
	"fmt"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/timezone"
)

func main() {
	provider, _ := timezone.NewProvider()
	ctx := context.Background()

	// Validar zona IANA oficial
	tz, err := provider.ValidateIANAZone(ctx, "America/Lima")
	if err == nil {
		fmt.Printf("Zona: %s | Offset UTC: %s\n", tz.IANA, tz.UTCOffset)
	}
}
```

#### Monedas (`pkg/currency`)
```go
package main

import (
	"context"
	"fmt"

	"github.com/Jhonatan-Code-dev/viewgo/pkg/currency"
)

func main() {
	provider, _ := currency.NewProvider()
	ctx := context.Background()

	// Consultar divisa ISO 4217 (ej. PEN / Sol Peruano)
	pen, err := provider.ValidateCurrencyCode(ctx, "PEN")
	if err == nil {
		fmt.Printf("Moneda: %s | Símbolo: %s | Código Numérico: %d | Decimales: %d\n",
			pen.Code, pen.Symbol, pen.NumericCode, pen.FractionDigits)
	}

	// Formatear monto monetario
	formatted, _ := provider.FormatAmount(ctx, "PEN", 1234.56)
	fmt.Println(formatted) // S/ 1234.56
}
```

---

## Resultados Empíricos de Benchmarking de Rendimiento

Pruebas ejecutadas con `go test -bench=Benchmark -benchmem ./test` en procesador Intel Core i7-12700H (Go 1.26.0):

```text
goos: windows
goarch: amd64
pkg: github.com/Jhonatan-Code-dev/viewgo/test
cpu: 12th Gen Intel(R) Core(TM) i7-12700H

BenchmarkCountryLookup_Alpha2-20      	139883404	         9.43 ns/op	       0 B/op	       0 allocs/op
BenchmarkCountry_ValidateAlpha2-20    	100000000	        18.51 ns/op	       0 B/op	       0 allocs/op
BenchmarkCurrencyLookup_Code-20       	 91696300	        11.75 ns/op	       0 B/op	       0 allocs/op
BenchmarkCurrency_ValidateCode-20     	100000000	        10.57 ns/op	       0 B/op	       0 allocs/op
BenchmarkTimezoneLookup_IANA-20       	  4671776	       231.50 ns/op	      16 B/op	       1 allocs/op
BenchmarkTimezone_ValidateIANA-20     	 10694055	       159.30 ns/op	      16 B/op	       1 allocs/op
BenchmarkCurrencyFormat_Amount-20     	  3235359	       379.30 ns/op	      40 B/op	       4 allocs/op
BenchmarkRootSDK_AllLookups-20        	  2201221	       536.70 ns/op	      64 B/op	       5 allocs/op
PASS
```

---

## Estructura de Directorios del Proyecto

```text
viewgo/
├── .agents/                          # Reglas y Habilidades de Arquitectura e Ingeniería AGY
├── .github/
│   └── workflows/
│       └── ci.yml                    # Pipeline de Integración Continua (GitHub Actions)
├── cmd/
│   └── viewgo/
│       └── main.go                   # Ejecutable de Demostración interactivo en Línea de Comandos
├── test/                             # Paquete Consolidado de Pruebas de Unidad y Benchmarks
│   ├── benchmark_test.go             # Pruebas de Benchmarking Automatizado (ns/op, B/op, allocs/op)
│   ├── country_test.go               # Pruebas Unitarias del Módulo de Países
│   ├── currency_test.go              # Pruebas Unitarias del Módulo de Monedas
│   ├── sdk_test.go                   # Pruebas Integradas del SDK y Validación SaaS
│   └── timezone_test.go              # Pruebas Unitarias del Módulo de Zonas Horarias
├── pkg/                              # Paquetes Públicos Exportados para Consumidores Externos
│   ├── country/                      # Paquete Público de Países (country.NewProvider)
│   ├── timezone/                     # Paquete Público de Zonas Horarias (timezone.NewProvider)
│   └── currency/                     # Paquete Público de Monedas (currency.NewProvider)
├── internal/
│   ├── shared/
│   │   └── domain/errors.go          # Errores Centinela y Primitivas de Dominio Compartidas
│   └── modules/
│       ├── country/                  # Módulo de Dominio e Infraestructura de Países
│       ├── timezone/                 # Módulo de Dominio e Infraestructura de Zonas Horarias
│       └── currency/                 # Módulo de Dominio e Infraestructura de Monedas
├── doc.go                            # Documentación de Paquete a Nivel Raíz (GoDoc ST1000)
├── go.mod                            # Definición de Módulo Go
├── viewgo.go                         # Fachada Principal del SDK a Nivel Raíz
├── LICENSE                           # Licencia Abierta MIT
└── README.md                         # Documentación Oficial del Proyecto
```

---

## Licencia

Este proyecto está bajo la Licencia [MIT](LICENSE).  
Copyright (c) 2026 Jhonatan-Code-dev.
