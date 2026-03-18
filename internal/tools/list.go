package tools

import (
	"strings"
)

// ListTimezones returns a filtered list of IANA timezone identifiers.
func ListTimezones(region, prefix string, limit int) ([]string, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	// Standard IANA timezone list (subset - can be extended)
	timezones := []string{
		"UTC",
		"Africa/Cairo", "Africa/Lagos", "Africa/Johannesburg", "Africa/Nairobi",
		"America/New_York", "America/Chicago", "America/Denver", "America/Los_Angeles",
		"America/Mexico_City", "America/Toronto", "America/Vancouver", "America/Sao_Paulo",
		"America/Buenos_Aires", "America/Anchorage", "America/Juneau",
		"Asia/Tokyo", "Asia/Shanghai", "Asia/Hong_Kong", "Asia/Singapore",
		"Asia/Bangkok", "Asia/Jakarta", "Asia/Manila", "Asia/Seoul",
		"Asia/Mumbai", "Asia/Kolkata", "Asia/Karachi", "Asia/Almaty",
		"Asia/Dubai", "Asia/Beirut", "Asia/Riyadh", "Asia/Tehran",
		"Australia/Sydney", "Australia/Melbourne", "Australia/Brisbane",
		"Australia/Perth", "Australia/Adelaide", "Pacific/Auckland",
		"Europe/London", "Europe/Paris", "Europe/Berlin", "Europe/Rome",
		"Europe/Amsterdam", "Europe/Madrid", "Europe/Vienna", "Europe/Prague",
		"Europe/Moscow", "Europe/Istanbul", "Europe/Athens", "Pacific/Fiji",
	}

	var filtered []string

	for _, tz := range timezones {
		// Apply region filter
		if region != "" {
			if !strings.HasPrefix(tz, region) {
				continue
			}
		}

		// Apply prefix filter (case-insensitive substring match)
		if prefix != "" {
			if !strings.Contains(strings.ToLower(tz), strings.ToLower(prefix)) {
				continue
			}
		}

		filtered = append(filtered, tz)

		if len(filtered) >= limit {
			break
		}
	}

	return filtered, nil
}
