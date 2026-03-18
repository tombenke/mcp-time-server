package tools

import (
	"fmt"
	"time"
)

// TimezoneInfo returns information about a timezone at a specific moment in time.
func TimezoneInfo(timezone, timestamp string) (map[string]any, error) {
	if timezone == "" {
		timezone = "UTC"
	}

	if err := ValidateTimezone(timezone); err != nil {
		return nil, err
	}

	var t time.Time

	if timestamp == "" {
		t = time.Now()
	} else {
		if err := ValidateRFC3339(timestamp); err != nil {
			return nil, fmt.Errorf("invalid timestamp: %w", err)
		}
		var parseErr error
		t, parseErr = time.Parse(time.RFC3339, timestamp)
		if parseErr != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", parseErr)
		}
	}

	loc, _ := time.LoadLocation(timezone)
	t = t.In(loc)

	name, offsetSecs := t.Zone()

	// Determine DST state by comparing with previous and next time
	// Check if DST is active by comparing offset with a non-DST reference
	testTime := time.Date(t.Year(), 1, 1, 12, 0, 0, 0, loc)
	_, winterOffset := testTime.Zone()
	_, currentOffset := t.Zone()
	dstActive := currentOffset != winterOffset

	hours := offsetSecs / 3600
	minutes := (offsetSecs % 3600) / 60
	offsetStr := fmt.Sprintf("%+03d:%02d", hours, minutes)

	result := map[string]any{
		"utc_offset_seconds": offsetSecs,
		"utc_offset_string":  offsetStr,
		"dst_active":         dstActive,
		"abbreviation":       name,
		"timezone":           timezone,
	}

	return result, nil
}
