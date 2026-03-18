package tools

import (
	"fmt"
	"time"
)

// Convert converts a timestamp from one timezone to another.
func Convert(timestamp, fromTimezone, toTimezone, format string) (map[string]any, error) {
	if fromTimezone == "" {
		fromTimezone = "UTC"
	}
	if toTimezone == "" {
		toTimezone = "UTC"
	}

	if err := ValidateRFC3339(timestamp); err != nil {
		return nil, err
	}
	if err := ValidateTimezone(fromTimezone); err != nil {
		return nil, fmt.Errorf("invalid from_timezone: %w", err)
	}
	if err := ValidateTimezone(toTimezone); err != nil {
		return nil, fmt.Errorf("invalid to_timezone: %w", err)
	}

	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse timestamp: %w", err)
	}

	fromLoc, _ := time.LoadLocation(fromTimezone)
	toLoc, _ := time.LoadLocation(toTimezone)

	// Parse the timestamp in the source timezone context
	src := t.In(fromLoc)

	// Get offset difference
	_, fromOffset := src.Zone()
	_, toOffset := src.In(toLoc).Zone()
	offsetDelta := toOffset - fromOffset

	// Convert to destination timezone
	converted := src.In(toLoc)

	result := map[string]any{
		"converted_time":       converted.Format(time.RFC3339),
		"source_timezone":      fromTimezone,
		"destination_timezone": toTimezone,
		"offset_delta_seconds": offsetDelta,
	}

	return result, nil
}
