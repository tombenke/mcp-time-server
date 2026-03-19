package tools

import (
	"testing"
)

func TestNow(t *testing.T) {
	tests := []struct {
		name      string
		tz        string
		format    string
		wantError bool
		fields    []string
	}{
		{
			name:      "UTC default",
			tz:        "",
			format:    "",
			wantError: false,
			fields:    []string{"time_rfc3339", "time_unix", "timezone"},
		},
		{
			name:      "Explicit US Eastern",
			tz:        "America/New_York",
			format:    "",
			wantError: false,
			fields:    []string{"time_rfc3339", "time_unix", "timezone"},
		},
		{
			name:      "Invalid timezone",
			tz:        "Invalid/Zone",
			format:    "",
			wantError: true,
			fields:    []string{},
		},
		{
			name:      "Tokyo timezone",
			tz:        "Asia/Tokyo",
			format:    "",
			wantError: false,
			fields:    []string{"time_rfc3339", "time_unix", "timezone"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Now(tt.tz, tt.format)

			if (err != nil) != tt.wantError {
				t.Errorf("Now() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				for _, field := range tt.fields {
					if _, ok := result[field]; !ok {
						t.Errorf("Now() missing field %q in result", field)
					}
				}

				expected := tt.tz
				if expected == "" {
					expected = "UTC"
				}
				if result["timezone"] != expected {
					t.Errorf("Now() timezone = %v, want %v", result["timezone"], tt.tz)
				}
			}
		})
	}
}
