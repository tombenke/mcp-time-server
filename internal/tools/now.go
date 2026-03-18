package tools

import (
	"time"
)

// Now returns the current time in the specified timezone.
func Now(timezone, format string) (map[string]any, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	if err := ValidateTimezone(timezone); err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)

	result := map[string]any{
		"time_rfc3339": now.Format(time.RFC3339),
		"time_unix":    now.Unix(),
		"timezone":     timezone,
	}

	return result, nil
}
