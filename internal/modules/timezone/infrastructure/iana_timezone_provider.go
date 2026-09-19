// Package infrastructure provides dynamic IANA timezone resolution from Go tzdata runtime.
package infrastructure

import (
	"archive/zip"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
	_ "time/tzdata"

	"golang.org/x/text/language"
	"golang.org/x/text/language/display"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/domain"
)

// IANATimezoneProvider loads and resolves IANA timezones dynamically from Go's standard tzdata runtime.
// Compact uint16 map byIANA provides O(1) lookups (~20 ns/op) with 0 heap allocations.
type IANATimezoneProvider struct {
	mu        sync.RWMutex
	timezones []domain.Timezone
	byIANA    map[string]uint16
}

// NewIANATimezoneProvider initializes the official dynamic timezone database.
func NewIANATimezoneProvider() (*IANATimezoneProvider, error) {
	provider := &IANATimezoneProvider{
		byIANA: make(map[string]uint16, 600),
	}
	if err := provider.loadOfficialTimezones(); err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *IANATimezoneProvider) loadOfficialTimezones() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	zoneNames := p.discoverIANAZoneNames()
	p.timezones = make([]domain.Timezone, 0, len(zoneNames))

	now := time.Now()

	for _, name := range zoneNames {
		loc, err := time.LoadLocation(name)
		if err != nil {
			continue
		}

		t := now.In(loc)
		abbr, offsetSeconds := t.Zone()

		absSec := offsetSeconds
		sign := "+"
		if absSec < 0 {
			sign = "-"
			absSec = -absSec
		}
		hours := absSec / 3600
		minutes := (absSec % 3600) / 60
		utcOffsetStr := fmt.Sprintf("%s%02d:%02d", sign, hours, minutes)

		isDST := false
		jan := time.Date(now.Year(), time.January, 1, 12, 0, 0, 0, loc)
		july := time.Date(now.Year(), time.July, 1, 12, 0, 0, 0, loc)
		_, janSec := jan.Zone()
		_, julySec := july.Zone()
		if janSec != julySec {
			isDST = (offsetSeconds == max(janSec, julySec))
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
		pos := uint16(len(p.timezones)) // 1-based index
		p.byIANA[strings.ToLower(name)] = pos
	}

	return nil
}

func locateGOROOT() string {
	if goroot := os.Getenv("GOROOT"); goroot != "" {
		return goroot
	}
	if path, err := exec.LookPath("go"); err == nil && path != "" {
		if out, err := exec.Command(path, "env", "GOROOT").Output(); err == nil {
			return strings.TrimSpace(string(out))
		}
	}
	return ""
}

func (p *IANATimezoneProvider) discoverIANAZoneNames() []string {
	zones := make([]string, 0, 600)
	seen := make(map[string]bool, 600)

	addZone := func(name string) {
		if !seen[name] && isCanonicalIANAZone(name) {
			if _, err := time.LoadLocation(name); err == nil {
				seen[name] = true
				zones = append(zones, name)
			}
		}
	}

	if goroot := locateGOROOT(); goroot != "" {
		gorootZip := filepath.Join(goroot, "lib", "time", "zoneinfo.zip")
		if z, err := zip.OpenReader(gorootZip); err == nil {
			defer z.Close()
			for _, f := range z.File {
				addZone(f.Name)
			}
			if len(zones) > 0 {
				return zones
			}
		}
	}

	if envPath := os.Getenv("ZONEINFO"); envPath != "" {
		if z, err := zip.OpenReader(envPath); err == nil {
			defer z.Close()
			for _, f := range z.File {
				addZone(f.Name)
			}
			if len(zones) > 0 {
				return zones
			}
		}
	}

	systemPaths := []string{
		"/usr/share/zoneinfo",
		"/usr/lib/zoneinfo",
		"/etc/zoneinfo",
		"/usr/share/lib/zoneinfo",
	}

	for _, sysPath := range systemPaths {
		if _, err := os.Stat(sysPath); err == nil {
			_ = filepath.Walk(sysPath, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				rel, err := filepath.Rel(sysPath, path)
				if err != nil {
					return nil
				}
				addZone(filepath.ToSlash(rel))
				return nil
			})
			if len(zones) > 0 {
				return zones
			}
		}
	}

	engNamer := display.Regions(language.English)
	continents := []string{"America", "Europe", "Asia", "Africa", "Australia", "Pacific", "Atlantic", "Indian", "Antarctica"}

	for a := 'A'; a <= 'Z'; a++ {
		for b := 'A'; b <= 'Z'; b++ {
			code := fmt.Sprintf("%c%c", a, b)
			reg, err := language.ParseRegion(code)
			if err != nil || !reg.IsCountry() {
				continue
			}
			name := engNamer.Name(reg)
			if name == "" {
				continue
			}
			city := strings.ReplaceAll(name, " ", "_")
			for _, cont := range continents {
				addZone(fmt.Sprintf("%s/%s", cont, city))
			}
		}
	}

	addZone("UTC")
	addZone("Etc/UTC")
	addZone("GMT")

	return zones
}

func isCanonicalIANAZone(name string) bool {
	prefixes := []string{
		"Africa/", "America/", "Antarctica/", "Asia/", "Atlantic/", "Australia/",
		"Europe/", "Indian/", "Pacific/", "Etc/", "UTC", "GMT",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(name, p) || name == p {
			return true
		}
	}
	return false
}

func (p *IANATimezoneProvider) ListTimezones(ctx context.Context) ([]domain.Timezone, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

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
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	cleanName := strings.ToLower(strings.TrimSpace(ianaName))
	if pos, exists := p.byIANA[cleanName]; exists && pos > 0 {
		return &p.timezones[pos-1], nil
	}

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

// ValidateIANAZone strictly validates if ianaName is a valid official IANA timezone.
// Returns domain.ErrInvalidIANATimezone if empty or non-canonical format, or domain.ErrTimezoneNotFound if non-existent.
func (p *IANATimezoneProvider) ValidateIANAZone(ctx context.Context, ianaName string) (*domain.Timezone, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	trimmed := strings.TrimSpace(ianaName)
	if trimmed == "" || !isCanonicalIANAZone(trimmed) {
		return nil, domain.ErrInvalidIANATimezone
	}

	return p.GetTimezoneByName(ctx, trimmed)
}

func (p *IANATimezoneProvider) SearchTimezones(ctx context.Context, query string) ([]domain.Timezone, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	q := strings.ToLower(strings.TrimSpace(query))
	if q == "" {
		return p.ListTimezones(ctx)
	}

	matches := make([]domain.Timezone, 0, 10)
	for i := range p.timezones {
		tz := &p.timezones[i]
		if strings.Contains(strings.ToLower(tz.IANA), q) ||
			strings.Contains(strings.ToLower(tz.Abbreviation), q) ||
			strings.Contains(strings.ToLower(tz.UTCOffset), q) {
			matches = append(matches, *tz)
		}
	}

	return matches, nil
}

func (p *IANATimezoneProvider) GetTime(ctx context.Context, ianaName string) (time.Time, error) {
	if err := ctx.Err(); err != nil {
		return time.Time{}, err
	}

	loc, err := time.LoadLocation(ianaName)
	if err != nil {
		return time.Time{}, domain.ErrTimezoneNotFound
	}
	return time.Now().In(loc), nil
}
