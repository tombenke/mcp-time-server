package tools

import (
	"fmt"
	"time"
)

// ValidateTimezone checks if the given timezone string is a valid IANA timezone.
func ValidateTimezone(tz string) error {
	if tz == "" {
		return nil // empty is treated as UTC by default
	}
	_, err := time.LoadLocation(tz)
	if err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}
	return nil
}

// ValidateRFC3339 checks if the given string is a valid RFC3339 timestamp.
func ValidateRFC3339(ts string) error {
	if ts == "" {
		return fmt.Errorf("timestamp cannot be empty")
	}
	_, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return fmt.Errorf("invalid RFC3339 timestamp: %w", err)
	}
	return nil
}
