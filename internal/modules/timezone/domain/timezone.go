// Package domain defines core business entities and provider interfaces for timezones.
package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrTimezoneNotFound = errors.New("timezone not found in IANA database")
)

// Timezone represents an official IANA time zone (Olson database) with dynamic offset details.
type Timezone struct {
	IANA             string `json:"iana"`               // Official IANA Zone identifier (e.g. "America/New_York", "Europe/Paris")
	Abbreviation     string `json:"abbreviation"`       // Dynamic zone abbreviation (e.g. "EST", "EDT", "CET", "JST")
	UTCOffset        string `json:"utc_offset"`         // Formatted UTC offset string (e.g. "-05:00", "+01:00", "+09:00")
	RawOffsetSeconds int    `json:"raw_offset_seconds"` // Current offset in seconds from UTC
	IsDST            bool   `json:"is_dst"`             // Whether Daylight Saving Time is currently active
	CurrentTime      string `json:"current_time"`       // ISO 8601 current timestamp in this timezone
}

// TimezoneProvider defines the repository interface for querying dynamic IANA timezones.
type TimezoneProvider interface {
	ListTimezones(ctx context.Context) ([]Timezone, error)
	GetTimezoneByName(ctx context.Context, ianaName string) (*Timezone, error)
	SearchTimezones(ctx context.Context, query string) ([]Timezone, error)
	GetTime(ctx context.Context, ianaName string) (time.Time, error)
}
