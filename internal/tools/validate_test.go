package tools

import (
	"testing"
)

func TestValidateTimezone(t *testing.T) {
	tests := []struct {
		name      string
		tz        string
		wantError bool
	}{
		{"UTC valid", "UTC", false},
		{"US Eastern valid", "America/New_York", false},
		{"Asia Tokyo valid", "Asia/Tokyo", false},
		{"Empty string (defaults)", "", false},
		{"Invalid timezone", "Invalid/Zone", true},
		{"Typo in timezone", "America/New_Yor", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTimezone(tt.tz)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateTimezone() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateRFC3339(t *testing.T) {
	tests := []struct {
		name      string
		ts        string
		wantError bool
	}{
		{"Valid RFC3339", "2024-03-17T10:30:00Z", false},
		{"Valid RFC3339 with offset", "2024-03-17T10:30:00+05:30", false},
		{"Invalid format", "2024-03-17 10:30:00", true},
		{"Invalid date", "2024-13-32T10:30:00Z", true},
		{"Empty string", "", true},
		{"Partial timestamp", "2024-03-17", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRFC3339(tt.ts)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateRFC3339() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
