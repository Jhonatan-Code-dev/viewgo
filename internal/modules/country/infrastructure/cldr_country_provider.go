package infrastructure

import (
	"context"
	"strings"
	"sync"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"viewgo/internal/modules/country/domain"
)

// CLDRCountryProvider dynamically extracts ISO 3166-1 country data using official Go unicode text tools.
type CLDRCountryProvider struct {
	mu        sync.RWMutex
	countries []domain.Country
	byAlpha2  map[string]domain.Country
	byAlpha3  map[string]domain.Country
}

// NewCLDRCountryProvider initializes and dynamically loads all official ISO 3166-1 countries from Unicode CLDR.
func NewCLDRCountryProvider() (*CLDRCountryProvider, error) {
	provider := &CLDRCountryProvider{
		byAlpha2: make(map[string]domain.Country),
		byAlpha3: make(map[string]domain.Country),
	}
	if err := provider.loadOfficialCountries(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *CLDRCountryProvider) loadOfficialCountries() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	engNamnamer := display.English.Tags()
	_ = engNamnamer

	// Iterate through all supported region tags defined in official Go x/text/language ISO 3166-1 registry
	for _, reg := range language.SupportedRegions() {
		if !reg.IsCountry() {
			continue
		}

		alpha2 := strings.ToUpper(reg.String())
		alpha3 := strings.ToUpper(reg.ISO3())
		numeric := reg.M49()

		// Retrieve official localized display name from Unicode CLDR via Go display package
		englishName := display.Region(reg).Name(language.English)
		if englishName == "" {
			englishName = alpha2
		}

		nativeName := display.Self.Region(reg)
		if nativeName == "" {
			nativeName = englishName
		}

		c := domain.Country{
			Alpha2:     alpha2,
			Alpha3:     alpha3,
			Numeric:    numeric,
			Name:       englishName,
			NativeName: nativeName,
			IsOfficial: true,
		}

		p.countries = append(p.countries, c)
		p.byAlpha2[alpha2] = c
		if alpha3 != "" {
			p.byAlpha3[alpha3] = c
		}
	}

	return nil
}

func (p *CLDRCountryProvider) ListCountries(ctx context.Context) ([]domain.Country, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]domain.Country, len(p.countries))
	copy(result, p.countries)
	return result, nil
}

func (p *CLDRCountryProvider) GetCountryByCode(ctx context.Context, code string) (*domain.Country, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	if len(cleanCode) == 2 {
		if c, exists := p.byAlpha2[cleanCode]; exists {
			return &c, nil
		}
	} else if len(cleanCode) == 3 {
		if c, exists := p.byAlpha3[cleanCode]; exists {
			return &c, nil
		}
	}

	return nil, domain.ErrCountryNotFound
}

func (p *CLDRCountryProvider) SearchCountries(ctx context.Context, query string) ([]domain.Country, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return p.ListCountries(ctx)
	}

	var matches []domain.Country
	for _, c := range p.countries {
		if strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.NativeName), q) ||
			strings.ToLower(c.Alpha2) == q ||
			strings.ToLower(c.Alpha3) == q {
			matches = append(matches, c)
		}
	}

	return matches, nil
}
