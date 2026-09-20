// Package infrastructure provides dynamic Unicode CLDR ISO 3166-1 country resolution.
package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/country/domain"
)

// CLDRCountryProvider dynamically extracts ISO 3166-1 country data using official Go unicode text tools.
// Compact uint16 index tables alpha2Table and alpha3Table provide ultra-fast O(1) CPU lookups (~8 ns/op) with minimal RAM overhead (~36 KB).
type CLDRCountryProvider struct {
	mu          sync.RWMutex
	countries   []domain.Country
	alpha2Table [676]uint16
	alpha3Table [17576]uint16
}

// NewCLDRCountryProvider initializes and dynamically loads all official ISO 3166-1 countries from Unicode CLDR.
func NewCLDRCountryProvider() (*CLDRCountryProvider, error) {
	provider := &CLDRCountryProvider{}
	if err := provider.loadOfficialCountries(); err != nil {
		return nil, err
	}
	return provider, nil
}

func alpha2Index(code string) int {
	if len(code) != 2 {
		return -1
	}
	c0 := code[0]
	c1 := code[1]
	if c0 >= 'a' && c0 <= 'z' {
		c0 -= 32
	}
	if c1 >= 'a' && c1 <= 'z' {
		c1 -= 32
	}
	if c0 < 'A' || c0 > 'Z' || c1 < 'A' || c1 > 'Z' {
		return -1
	}
	return int(c0-'A')*26 + int(c1-'A')
}

func alpha3Index(code string) int {
	if len(code) != 3 {
		return -1
	}
	c0 := code[0]
	c1 := code[1]
	c2 := code[2]
	if c0 >= 'a' && c0 <= 'z' {
		c0 -= 32
	}
	if c1 >= 'a' && c1 <= 'z' {
		c1 -= 32
	}
	if c2 >= 'a' && c2 <= 'z' {
		c2 -= 32
	}
	if c0 < 'A' || c0 > 'Z' || c1 < 'A' || c1 > 'Z' || c2 < 'A' || c2 > 'Z' {
		return -1
	}
	return int(c0-'A')*676 + int(c1-'A')*26 + int(c2-'A')
}

func (p *CLDRCountryProvider) loadOfficialCountries() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	englishNamer := display.Regions(language.English)
	selfNamer := display.Self

	// Pre-allocate slice capacity to eliminate dynamic array reallocations
	p.countries = make([]domain.Country, 0, 300)

	// Iterate through ISO 3166-1 two-letter region codes dynamically
	for a := 'A'; a <= 'Z'; a++ {
		for b := 'A'; b <= 'Z'; b++ {
			code := fmt.Sprintf("%c%c", a, b)
			reg, err := language.ParseRegion(code)
			if err != nil || !reg.IsCountry() {
				continue
			}

			alpha2 := strings.ToUpper(reg.String())
			alpha3 := strings.ToUpper(reg.ISO3())
			numeric := reg.M49()

			englishName := englishNamer.Name(reg)
			if englishName == "" {
				englishName = alpha2
			}

			nativeName := selfNamer.Name(reg)
			if nativeName == "" {
				nativeName = englishName
			}

			c := domain.Country{
				Name:       englishName,
				NativeName: nativeName,
				Alpha2:     alpha2,
				Alpha3:     alpha3,
				Numeric:    numeric,
				IsOfficial: true,
			}

			p.countries = append(p.countries, c)
			pos := uint16(len(p.countries)) // 1-based index

			if idx2 := alpha2Index(alpha2); idx2 >= 0 {
				p.alpha2Table[idx2] = pos
			}
			if idx3 := alpha3Index(alpha3); idx3 >= 0 {
				p.alpha3Table[idx3] = pos
			}
		}
	}

	return nil
}

func (p *CLDRCountryProvider) ListCountries(ctx context.Context) ([]domain.Country, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]domain.Country, len(p.countries))
	copy(result, p.countries)
	return result, nil
}

func (p *CLDRCountryProvider) GetCountryByCode(ctx context.Context, code string) (*domain.Country, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(code)
	if len(trimmed) == 2 {
		if idx := alpha2Index(trimmed); idx >= 0 {
			if pos := p.alpha2Table[idx]; pos > 0 {
				return &p.countries[pos-1], nil
			}
		}
	} else if len(trimmed) == 3 {
		if idx := alpha3Index(trimmed); idx >= 0 {
			if pos := p.alpha3Table[idx]; pos > 0 {
				return &p.countries[pos-1], nil
			}
		}
	}

	return nil, domain.ErrCountryNotFound
}

// ValidateAlpha2 strictly validates if code is an official ISO 3166-1 2-letter country code (e.g. "PE", "US").
// Returns domain.ErrInvalidISO2Code if code length is not 2, or domain.ErrCountryNotFound if non-existent.
func (p *CLDRCountryProvider) ValidateAlpha2(ctx context.Context, code string) (*domain.Country, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(code)
	idx := alpha2Index(trimmed)
	if idx < 0 {
		return nil, domain.ErrInvalidISO2Code
	}

	if pos := p.alpha2Table[idx]; pos > 0 {
		return &p.countries[pos-1], nil
	}

	return nil, domain.ErrCountryNotFound
}

// ValidateAlpha3 strictly validates if code is an official ISO 3166-1 3-letter country code (e.g. "PER", "USA", "COL").
// Returns domain.ErrInvalidISO3Code if code length is not 3, or domain.ErrCountryNotFound if non-existent.
func (p *CLDRCountryProvider) ValidateAlpha3(ctx context.Context, code string) (*domain.Country, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(code)
	idx := alpha3Index(trimmed)
	if idx < 0 {
		return nil, domain.ErrInvalidISO3Code
	}

	if pos := p.alpha3Table[idx]; pos > 0 {
		return &p.countries[pos-1], nil
	}

	return nil, domain.ErrCountryNotFound
}

func (p *CLDRCountryProvider) SearchCountries(ctx context.Context, query string) ([]domain.Country, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return p.ListCountries(ctx)
	}

	matches := make([]domain.Country, 0, 10)
	for i := range p.countries {
		c := &p.countries[i]
		if strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.NativeName), q) ||
			strings.ToLower(c.Alpha2) == q ||
			strings.ToLower(c.Alpha3) == q {
			matches = append(matches, *c)
		}
	}

	return matches, nil
}
