package tools

import (
	"testing"
)

func TestTimezoneInfo(t *testing.T) {
	tests := []struct {
		name       string
		tz         string
		timestamp  string
		wantError  bool
		wantFields []string
	}{
		{
			name:       "UTC zone",
			tz:         "UTC",
			timestamp:  "",
			wantError:  false,
			wantFields: []string{"utc_offset_seconds", "utc_offset_string", "dst_active", "abbreviation", "timezone"},
		},
		{
			name:       "US Eastern zone",
			tz:         "America/New_York",
			timestamp:  "",
			wantError:  false,
			wantFields: []string{"utc_offset_seconds", "utc_offset_string", "dst_active", "abbreviation", "timezone"},
		},
		{
			name:       "With explicit timestamp",
			tz:         "America/New_York",
			timestamp:  "2024-06-17T12:00:00Z", // Summer (DST active)
			wantError:  false,
			wantFields: []string{"utc_offset_seconds", "utc_offset_string", "dst_active", "abbreviation", "timezone"},
		},
		{
			name:       "Invalid timezone",
			tz:         "Invalid/Zone",
			timestamp:  "",
			wantError:  true,
			wantFields: []string{},
		},
		{
			name:       "Invalid timestamp",
			tz:         "UTC",
			timestamp:  "not-a-timestamp",
			wantError:  true,
			wantFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := TimezoneInfo(tt.tz, tt.timestamp)

			if (err != nil) != tt.wantError {
				t.Errorf("TimezoneInfo() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				for _, field := range tt.wantFields {
					if _, ok := result[field]; !ok {
						t.Errorf("TimezoneInfo() missing field %q in result", field)
					}
				}
			}
		})
	}
}
