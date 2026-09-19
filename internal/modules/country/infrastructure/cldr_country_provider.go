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
// Pointer maps byAlpha2 and byAlpha3 provide sub-microsecond O(1) lookups with 0 heap allocations.
type CLDRCountryProvider struct {
	mu        sync.RWMutex
	countries []domain.Country
	byAlpha2  map[string]*domain.Country
	byAlpha3  map[string]*domain.Country
}

// NewCLDRCountryProvider initializes and dynamically loads all official ISO 3166-1 countries from Unicode CLDR.
func NewCLDRCountryProvider() (*CLDRCountryProvider, error) {
	provider := &CLDRCountryProvider{
		byAlpha2: make(map[string]*domain.Country, 300),
		byAlpha3: make(map[string]*domain.Country, 300),
	}
	if err := provider.loadOfficialCountries(); err != nil {
		return nil, err
	}
	return provider, nil
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
			ptr := &p.countries[len(p.countries)-1]
			p.byAlpha2[alpha2] = ptr
			if alpha3 != "" {
				p.byAlpha3[alpha3] = ptr
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

	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	if len(cleanCode) == 2 {
		if ptr, exists := p.byAlpha2[cleanCode]; exists {
			return ptr, nil
		}
	} else if len(cleanCode) == 3 {
		if ptr, exists := p.byAlpha3[cleanCode]; exists {
			return ptr, nil
		}
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
