package tools

import (
	"testing"
)

func TestConvert(t *testing.T) {
	tests := []struct {
		name          string
		timestamp     string
		fromTZ        string
		toTZ          string
		format        string
		wantError     bool
		wantFields    []string
	}{
		{
			name:        "Round-trip same zone",
			timestamp:   "2024-03-17T10:00:00Z",
			fromTZ:      "UTC",
			toTZ:        "UTC",
			format:      "",
			wantError:   false,
			wantFields:  []string{"converted_time", "source_timezone", "destination_timezone"},
		},
		{
			name:        "Cross-zone conversion",
			timestamp:   "2024-03-17T10:00:00Z",
			fromTZ:      "UTC",
			toTZ:        "America/New_York",
			format:      "",
			wantError:   false,
			wantFields:  []string{"converted_time", "source_timezone", "destination_timezone"},
		},
		{
			name:        "Invalid from-zone",
			timestamp:   "2024-03-17T10:00:00Z",
			fromTZ:      "Invalid/Zone",
			toTZ:        "UTC",
			format:      "",
			wantError:   true,
			wantFields:  []string{},
		},
		{
			name:        "Invalid to-zone",
			timestamp:   "2024-03-17T10:00:00Z",
			fromTZ:      "UTC",
			toTZ:        "Invalid/Zone",
			format:      "",
			wantError:   true,
			wantFields:  []string{},
		},
		{
			name:        "Malformed timestamp",
			timestamp:   "not-a-timestamp",
			fromTZ:      "UTC",
			toTZ:        "UTC",
			format:      "",
			wantError:   true,
			wantFields:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Convert(tt.timestamp, tt.fromTZ, tt.toTZ, tt.format)

			if (err != nil) != tt.wantError {
				t.Errorf("Convert() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				for _, field := range tt.wantFields {
					if _, ok := result[field]; !ok {
						t.Errorf("Convert() missing field %q in result", field)
					}
				}
			}
		})
	}
}
