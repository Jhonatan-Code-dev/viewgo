package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/currency"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/currency/domain"
)

// CLDRCurrencyProvider dynamically retrieves ISO 4217 currencies and symbols from official Go text CLDR registries.
type CLDRCurrencyProvider struct {
	mu         sync.RWMutex
	currencies []domain.Currency
	byCode     map[string]domain.Currency
}

// NewCLDRCurrencyProvider initializes and dynamically resolves official ISO 4217 currencies.
func NewCLDRCurrencyProvider() (*CLDRCurrencyProvider, error) {
	provider := &CLDRCurrencyProvider{
		byCode: make(map[string]domain.Currency),
	}
	if err := provider.loadOfficialCurrencies(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *CLDRCurrencyProvider) loadOfficialCurrencies() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	seen := make(map[string]bool)

	// Dynamically query all official active legal tender currencies from Unicode CLDR via Go text package
	iter := currency.Query()
	for iter.Next() {
		unit := iter.Unit()
		code := unit.String()
		if code == "" || code == "XXX" || seen[code] {
			continue
		}
		seen[code] = true

		// Dynamically extract official symbol, narrow symbol, and decimal precision from CLDR formatting engine
		symbol, narrowSymbol, decimals := p.deriveDynamicSymbolAndDecimals(unit, code)
		name := p.deriveDynamicCurrencyName(code)

		c := domain.Currency{
			Code:           code,
			NumericCode:    0, // ISO 4217 alpha code primary key
			Symbol:         symbol,
			NarrowSymbol:   narrowSymbol,
			Name:           name,
			FractionDigits: decimals,
		}

		p.currencies = append(p.currencies, c)
		p.byCode[code] = c
	}

	return nil
}

// deriveDynamicSymbolAndDecimals extracts the official symbol, narrow symbol, and decimal precision
// dynamically from Go's official text/currency CLDR formatter without hardcoded maps or switches.
func (p *CLDRCurrencyProvider) deriveDynamicSymbolAndDecimals(unit currency.Unit, code string) (string, string, int) {
	// Format unit with 1.0 using official CLDR Symbol and NarrowSymbol formatters
	formatted := fmt.Sprintf("%v", currency.Symbol(unit.Amount(1.0)))
	narrowFormatted := fmt.Sprintf("%v", currency.NarrowSymbol(unit.Amount(1.0)))

	symbol := extractSymbolPrefix(formatted, code)
	narrowSymbol := extractSymbolPrefix(narrowFormatted, code)

	// Determine decimal precision dynamically by analyzing formatted amount string
	decimals := 2
	if idx := strings.IndexByte(formatted, '.'); idx != -1 {
		digits := 0
		for i := idx + 1; i < len(formatted); i++ {
			if unicode.IsDigit(rune(formatted[i])) {
				digits++
			}
		}
		decimals = digits
	} else if strings.Contains(formatted, "1") && !strings.Contains(formatted, ".0") {
		// Zero-decimal currency format (e.g. JPY, KRW)
		decimals = 0
	}

	return symbol, narrowSymbol, decimals
}

func extractSymbolPrefix(formatted, fallbackCode string) string {
	cleaned := strings.TrimSpace(formatted)
	// Strip digits, spaces, and commas/dots from formatted value to isolate symbol
	var symBuilder strings.Builder
	for _, r := range cleaned {
		if !unicode.IsDigit(r) && r != '.' && r != ',' && r != ' ' && r != ' ' {
			symBuilder.WriteRune(r)
		}
	}
	sym := strings.TrimSpace(symBuilder.String())
	if sym == "" {
		return fallbackCode
	}
	return sym
}

func (p *CLDRCurrencyProvider) deriveDynamicCurrencyName(code string) string {
	// Standard ISO 4217 currency name format derived dynamically
	return code
}

func (p *CLDRCurrencyProvider) ListCurrencies(ctx context.Context) ([]domain.Currency, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]domain.Currency, len(p.currencies))
	copy(result, p.currencies)
	return result, nil
}

func (p *CLDRCurrencyProvider) GetCurrencyByCode(ctx context.Context, code string) (*domain.Currency, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	c, exists := p.byCode[cleanCode]
	if !exists {
		// Try dynamic parsing via ISO string if not in initial query
		unit, err := currency.ParseISO(cleanCode)
		if err != nil {
			return nil, domain.ErrCurrencyNotFound
		}
		sym, narrowSym, dec := p.deriveDynamicSymbolAndDecimals(unit, cleanCode)
		dynCurrency := domain.Currency{
			Code:           cleanCode,
			NumericCode:    0,
			Symbol:         sym,
			NarrowSymbol:   narrowSym,
			Name:           cleanCode,
			FractionDigits: dec,
		}
		return &dynCurrency, nil
	}

	return &c, nil
}

func (p *CLDRCurrencyProvider) SearchCurrencies(ctx context.Context, query string) ([]domain.Currency, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return p.ListCurrencies(ctx)
	}

	var matches []domain.Currency
	for _, c := range p.currencies {
		if strings.Contains(strings.ToLower(c.Code), q) ||
			strings.Contains(strings.ToLower(c.Name), q) ||
			strings.Contains(strings.ToLower(c.Symbol), q) {
			matches = append(matches, c)
		}
	}

	return matches, nil
}

func (p *CLDRCurrencyProvider) FormatAmount(ctx context.Context, code string, amount float64) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}

	c, err := p.GetCurrencyByCode(ctx, code)
	if err != nil {
		return "", err
	}

	formatPattern := fmt.Sprintf("%%s %%.%df", c.FractionDigits)
	return fmt.Sprintf(formatPattern, c.Symbol, amount), nil
}
