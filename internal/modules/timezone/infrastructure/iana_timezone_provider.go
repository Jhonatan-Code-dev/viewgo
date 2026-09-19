package infrastructure

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	_ "time/tzdata" // Embedded official IANA time zone database from Go standard runtime

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/domain"
)

// Standard list of official IANA canonical continent/ocean prefixes according to the IANA TZ database.
var canonicalIANAPrefixes = []string{
	"Africa/", "America/", "Antarctica/", "Asia/", "Atlantic/", "Australia/",
	"Europe/", "Indian/", "Pacific/", "UTC", "Etc/UTC", "GMT",
}

// IANATimezoneProvider loads and resolves IANA timezones dynamically from Go's standard tzdata runtime.
type IANATimezoneProvider struct {
	mu        sync.RWMutex
	timezones []domain.Timezone
	byIANA    map[string]domain.Timezone
}

// NewIANATimezoneProvider initializes the official dynamic timezone database.
func NewIANATimezoneProvider() (*IANATimezoneProvider, error) {
	provider := &IANATimezoneProvider{
		byIANA: make(map[string]domain.Timezone),
	}
	if err := provider.loadOfficialTimezones(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *IANATimezoneProvider) loadOfficialTimezones() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Comprehensive list of standard canonical IANA time zones across all continents and regions
	knownIANAZones := []string{
		"UTC", "Etc/UTC", "GMT",
		// America
		"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
		"America/Anchorage", "America/Adak", "America/Phoenix", "America/Toronto",
		"America/Vancouver", "America/Mexico_City", "America/Bogota", "America/Lima",
		"America/Santiago", "America/Buenos_Aires", "America/Sao_Paulo", "America/Caracas",
		"America/La_Paz", "America/Guayaquil", "America/Panama", "America/Costa_Rica",
		"America/Havana", "America/Montevideo", "America/Asuncion", "America/Puerto_Rico",
		"America/Jamaica", "America/Santo_Domingo", "America/Guatemala", "America/Tegucigalpa",
		"America/Managua", "America/El_Salvador", "America/Belize", "America/Halifax",
		"America/St_Johns", "America/Edmonton", "America/Winnipeg", "America/Regina",
		// Europe
		"Europe/London", "Europe/Dublin", "Europe/Lisbon", "Europe/Madrid",
		"Europe/Paris", "Europe/Brussels", "Europe/Amsterdam", "Europe/Berlin",
		"Europe/Rome", "Europe/Vienna", "Europe/Zurich", "Europe/Stockholm",
		"Europe/Oslo", "Europe/Copenhagen", "Europe/Helsinki", "Europe/Athens",
		"Europe/Istanbul", "Europe/Moscow", "Europe/Warsaw", "Europe/Prague",
		"Europe/Budapest", "Europe/Bucharest", "Europe/Kiev", "Europe/Belgrade",
		// Asia
		"Asia/Tokyo", "Asia/Seoul", "Asia/Shanghai", "Asia/Hong_Kong",
		"Asia/Singapore", "Asia/Bangkok", "Asia/Jakarta", "Asia/Manila",
		"Asia/Kuala_Lumpur", "Asia/Ho_Chi_Minh", "Asia/Kolkata", "Asia/Dhaka",
		"Asia/Karachi", "Asia/Dubai", "Asia/Riyadh", "Asia/Tehran",
		"Asia/Baghdad", "Asia/Jerusalem", "Asia/Beirut", "Asia/Amman",
		"Asia/Damascus", "Asia/Tashkent", "Asia/Almaty", "Asia/Taipei",
		// Australia & Pacific
		"Australia/Sydney", "Australia/Melbourne", "Australia/Brisbane", "Australia/Adelaide",
		"Australia/Perth", "Australia/Hobart", "Australia/Darwin", "Pacific/Auckland",
		"Pacific/Fiji", "Pacific/Honolulu", "Pacific/Guam", "Pacific/Port_Moresby",
		"Pacific/Tahiti", "Pacific/Samoa", "Pacific/Tongatapu",
		// Africa
		"Africa/Cairo", "Africa/Johannesburg", "Africa/Lagos", "Africa/Nairobi", "Africa/Casablanca",
		"Africa/Algiers", "Africa/Tunis", "Africa/Accra", "Africa/Addis_Ababa", "Africa/Khartoum",
		"Africa/Dakar", "Africa/Luanda", "Africa/Kinshasa", "Africa/Harare",
		// Atlantic & Indian & Antarctica
		"Atlantic/Azores", "Atlantic/Canary", "Atlantic/Bermuda", "Atlantic/Reykjavik",
		"Indian/Mauritius", "Indian/Maldives", "Indian/Madagascar", "Antarctica/Palmer",
	}

	now := time.Now()

	for _, name := range knownIANAZones {
		loc, err := time.LoadLocation(name)
		if err != nil {
			continue // Skip if not found in Go's standard tzdata
		}

		t := now.In(loc)
		abbr, offsetSeconds := t.Zone()

		// Calculate UTC offset string (+HH:MM or -HH:MM)
		absSec := offsetSeconds
		sign := "+"
		if absSec < 0 {
			sign = "-"
			absSec = -absSec
		}
		hours := absSec / 3600
		minutes := (absSec % 3600) / 60
		utcOffsetStr := fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)

		// Determine if DST is active by comparing with standard January/July offsets
		isDST := false
		jan := time.Date(now.Year(), time.January, 1, 12, 0, 0, 0, loc)
		july := time.Date(now.Year(), time.July, 1, 12, 0, 0, 0, loc)
		_, janSec := jan.Zone()
		_, julySec := july.Zone()
		if janSec != julySec {
			isDST = (offsetSeconds == max(janSec, julySec)) && (janSec != julySec)
		}

		tz := domain.Timezone{
			IANA:             name,
			Abbreviation:     abbr,
			UTCOffset:        utcOffsetStr,
			RawOffsetSeconds: offsetSeconds,
			IsDST:            isDST,
			CurrentTime:      t.Format(time.RFC3339),
		}

		p.timezones = append(p.timezones, tz)
		p.byIANA[strings.ToLower(name)] = tz
	}

	return nil
}

func (p *IANATimezoneProvider) ListTimezones(ctx context.Context) ([]domain.Timezone, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Update current timestamps dynamically
	now := time.Now()
	result := make([]domain.Timezone, len(p.timezones))
	for i, tz := range p.timezones {
		if loc, err := time.LoadLocation(tz.IANA); err == nil {
			t := now.In(loc)
			tz.CurrentTime = t.Format(time.RFC3339)
		}
		result[i] = tz
	}
	return result, nil
}

func (p *IANATimezoneProvider) GetTimezoneByName(ctx context.Context, ianaName string) (*domain.Timezone, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanName := strings.ToLower(strings.TrimSpace(ianaName))
	tz, exists := p.byIANA[cleanName]
	if !exists {
		// Try loading dynamically directly via time.LoadLocation
		loc, err := time.LoadLocation(ianaName)
		if err != nil {
			return nil, domain.ErrTimezoneNotFound
		}
		now := time.Now().In(loc)
		abbr, offsetSec := now.Zone()
		absSec := offsetSec
		sign := "+"
		if absSec < 0 {
			sign = "-"
			absSec = -absSec
		}
		utcOffsetStr := fmt.Sprintf("%s%02d:%02d", sign, absSec/3600, (absSec%3600)/60)
		dynamicTZ := domain.Timezone{
			IANA:             ianaName,
			Abbreviation:     abbr,
			UTCOffset:        utcOffsetStr,
			RawOffsetSeconds: offsetSec,
			IsDST:            false,
			CurrentTime:      now.Format(time.RFC3339),
		}
		return &dynamicTZ, nil
	}

	if loc, err := time.LoadLocation(tz.IANA); err == nil {
		t := time.Now().In(loc)
		tz.CurrentTime = t.Format(time.RFC3339)
	}

	return &tz, nil
}

func (p *IANATimezoneProvider) SearchTimezones(ctx context.Context, query string) ([]domain.Timezone, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return p.ListTimezones(ctx)
	}

	var matches []domain.Timezone
	for _, tz := range p.timezones {
		if strings.Contains(strings.ToLower(tz.IANA), q) ||
			strings.Contains(strings.ToLower(tz.Abbreviation), q) ||
			strings.Contains(strings.ToLower(tz.UTCOffset), q) {
			matches = append(matches, tz)
		}
	}

	return matches, nil
}

func (p *IANATimezoneProvider) GetTime(ctx context.Context, ianaName string) (time.Time, error) {
	loc, err := time.LoadLocation(ianaName)
	if err != nil {
		return time.Time{}, domain.ErrTimezoneNotFound
	}
	return time.Now().In(loc), nil
}
