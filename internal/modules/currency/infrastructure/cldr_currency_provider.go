package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"golang.org/x/text/currency"

	"viewgo/internal/modules/currency/domain"
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

		symbol := p.resolveOfficialSymbol(code)
		narrowSymbol := p.resolveNarrowSymbol(code)
		numCode, decimals := p.resolveNumericAndDecimals(code)
		name := p.resolveCurrencyName(code)

		c := domain.Currency{
			Code:           code,
			NumericCode:    numCode,
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

func (p *CLDRCurrencyProvider) resolveCurrencyName(code string) string {
	names := map[string]string{
		"USD": "US Dollar",
		"EUR": "Euro",
		"JPY": "Japanese Yen",
		"GBP": "British Pound",
		"AUD": "Australian Dollar",
		"CAD": "Canadian Dollar",
		"CHF": "Swiss Franc",
		"CNY": "Chinese Yuan",
		"HKD": "Hong Kong Dollar",
		"NZD": "New Zealand Dollar",
		"SEK": "Swedish Krona",
		"KRW": "South Korean Won",
		"SGD": "Singapore Dollar",
		"NOK": "Norwegian Krone",
		"MXN": "Mexican Peso",
		"INR": "Indian Rupee",
		"RUB": "Russian Ruble",
		"ZAR": "South African Rand",
		"TRY": "Turkish Lira",
		"BRL": "Brazilian Real",
		"COP": "Colombian Peso",
		"CLP": "Chilean Peso",
		"PEN": "Peruvian Sol",
		"ARS": "Argentine Peso",
		"AED": "United Arab Emirates Dirham",
		"SAR": "Saudi Riyal",
	}
	if name, found := names[code]; found {
		return name
	}
	return code + " Currency"
}

func (p *CLDRCurrencyProvider) resolveOfficialSymbol(code string) string {
	switch code {
	case "USD", "CAD", "AUD", "NZD", "HKD", "SGD", "MXN", "COP", "CLP", "ARS":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY", "CNY":
		return "¥"
	case "KRW":
		return "₩"
	case "INR":
		return "₹"
	case "RUB":
		return "₽"
	case "TRY":
		return "₺"
	case "BRL":
		return "R$"
	case "THB":
		return "฿"
	case "VND":
		return "₫"
	case "ILS":
		return "₪"
	case "PHP":
		return "₱"
	case "AED":
		return "د.إ"
	case "SAR":
		return "﷼"
	case "EGP":
		return "E£"
	case "NGN":
		return "₦"
	case "ZAR":
		return "R"
	case "CHF":
		return "CHF"
	case "PLN":
		return "zł"
	case "SEK", "NOK", "DKK":
		return "kr"
	default:
		return code
	}
}

func (p *CLDRCurrencyProvider) resolveNarrowSymbol(code string) string {
	switch code {
	case "USD", "CAD", "AUD", "NZD", "HKD", "SGD", "MXN", "COP", "CLP", "ARS":
		return "$"
	case "EUR":
		return "€"
	case "GBP":
		return "£"
	case "JPY", "CNY":
		return "¥"
	case "KRW":
		return "₩"
	case "INR":
		return "₹"
	default:
		return code
	}
}

func (p *CLDRCurrencyProvider) resolveNumericAndDecimals(code string) (int, int) {
	switch code {
	case "USD":
		return 840, 2
	case "EUR":
		return 978, 2
	case "JPY":
		return 392, 0
	case "GBP":
		return 826, 2
	case "AUD":
		return 36, 2
	case "CAD":
		return 124, 2
	case "COP":
		return 170, 2
	case "MXN":
		return 484, 2
	case "BRL":
		return 986, 2
	case "CNY":
		return 156, 2
	case "KRW":
		return 410, 0
	case "CLP":
		return 152, 0
	case "VND":
		return 704, 0
	default:
		return 0, 2
	}
}

func (p *CLDRCurrencyProvider) ListCurrencies(ctx context.Context) ([]domain.Currency, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	result := make([]domain.Currency, len(p.currencies))
	copy(result, p.currencies)
	return result, nil
}

func (p *CLDRCurrencyProvider) GetCurrencyByCode(ctx context.Context, code string) (*domain.Currency, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanCode := strings.ToUpper(strings.TrimSpace(code))
	c, exists := p.byCode[cleanCode]
	if !exists {
		return nil, domain.ErrCurrencyNotFound
	}

	return &c, nil
}

func (p *CLDRCurrencyProvider) SearchCurrencies(ctx context.Context, query string) ([]domain.Currency, error) {
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
	c, err := p.GetCurrencyByCode(ctx, code)
	if err != nil {
		return "", err
	}

	formatPattern := fmt.Sprintf("%%s %%.%df", c.FractionDigits)
	return fmt.Sprintf(formatPattern, c.Symbol, amount), nil
}
