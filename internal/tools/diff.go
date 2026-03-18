package tools

import (
	"fmt"
	"time"
)

// Diff calculates the time difference between two RFC3339 timestamps.
func Diff(startTimestamp, endTimestamp string) (map[string]any, error) {
	if err := ValidateRFC3339(startTimestamp); err != nil {
		return nil, fmt.Errorf("invalid start_timestamp: %w", err)
	}
	if err := ValidateRFC3339(endTimestamp); err != nil {
		return nil, fmt.Errorf("invalid end_timestamp: %w", err)
	}

	start, err := time.Parse(time.RFC3339, startTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse start_timestamp: %w", err)
	}

	end, err := time.Parse(time.RFC3339, endTimestamp)
	if err != nil {
		return nil, fmt.Errorf("failed to parse end_timestamp: %w", err)
	}

	duration := end.Sub(start)

	seconds := int64(duration.Seconds())
	minutes := duration.Minutes()
	hours := duration.Hours()

	absSeconds := seconds
	absMinutes := minutes
	absHours := hours

	if absSeconds < 0 {
		absSeconds = -absSeconds
	}
	if absMinutes < 0 {
		absMinutes = -absMinutes
	}
	if absHours < 0 {
		absHours = -absHours
	}

	result := map[string]any{
		"seconds":      seconds,
		"minutes":      minutes,
		"hours":        hours,
		"abs_seconds":  absSeconds,
		"abs_minutes":  absMinutes,
		"abs_hours":    absHours,
	}

	return result, nil
}
