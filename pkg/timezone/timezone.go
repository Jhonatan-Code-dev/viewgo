// Package timezone provides dynamic, official IANA timezone resolution derived from Go's tzdata runtime.
package timezone

import (
	"context"
	"time"

	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/domain"
	"github.com/Jhonatan-Code-dev/viewgo/internal/modules/timezone/infrastructure"
)

// Timezone represents an official IANA time zone identifier and offset metadata.
type Timezone = domain.Timezone

// Provider defines the interface for querying dynamic IANA timezones.
type Provider interface {
	ListTimezones(ctx context.Context) ([]Timezone, error)
	GetTimezoneByName(ctx context.Context, ianaName string) (*Timezone, error)
	ValidateIANAZone(ctx context.Context, ianaName string) (*Timezone, error)
	SearchTimezones(ctx context.Context, query string) ([]Timezone, error)
	GetTime(ctx context.Context, ianaName string) (time.Time, error)
}

// NewProvider initializes and returns a new official IANA Timezone Provider.
func NewProvider() (Provider, error) {
	return infrastructure.NewIANATimezoneProvider()
}
