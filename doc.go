// Package viewgo provides production-grade, dynamic, official reference data for
// Countries (ISO 3166-1), IANA Timezones, and Currencies & Symbols (ISO 4217).
//
// All reference data is dynamically loaded from official Go runtime packages (time/tzdata)
// and official Unicode CLDR repositories (golang.org/x/text) without hardcoded static datasets.
//
// Subpackages:
//   - pkg/country: ISO 3166-1 Country registry and localized names.
//   - pkg/timezone: IANA Timezone registry, UTC offsets, and DST calculations.
//   - pkg/currency: ISO 4217 Currency registry, symbols, and formatting tools.
package viewgo
