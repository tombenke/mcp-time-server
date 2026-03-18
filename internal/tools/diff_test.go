package tools

import (
	"testing"
)

func TestDiff(t *testing.T) {
	tests := []struct {
		name         string
		startTS      string
		endTS        string
		wantError    bool
		wantFields   []string
		validateSign bool // If true, check that seconds field has expected sign
		wantPositive bool // Expected sign if validateSign is true
	}{
		{
			name:         "Positive diff",
			startTS:      "2024-03-17T10:00:00Z",
			endTS:        "2024-03-17T11:00:00Z",
			wantError:    false,
			wantFields:   []string{"seconds", "minutes", "hours", "abs_seconds", "abs_minutes", "abs_hours"},
			validateSign: true,
			wantPositive: true,
		},
		{
			name:         "Negative diff",
			startTS:      "2024-03-17T11:00:00Z",
			endTS:        "2024-03-17T10:00:00Z",
			wantError:    false,
			wantFields:   []string{"seconds", "minutes", "hours", "abs_seconds", "abs_minutes", "abs_hours"},
			validateSign: true,
			wantPositive: false,
		},
		{
			name:       "Zero diff",
			startTS:    "2024-03-17T10:00:00Z",
			endTS:      "2024-03-17T10:00:00Z",
			wantError:  false,
			wantFields: []string{"seconds", "minutes", "hours", "abs_seconds", "abs_minutes", "abs_hours"},
		},
		{
			name:       "Malformed start timestamp",
			startTS:    "not-a-timestamp",
			endTS:      "2024-03-17T10:00:00Z",
			wantError:  true,
			wantFields: []string{},
		},
		{
			name:       "Malformed end timestamp",
			startTS:    "2024-03-17T10:00:00Z",
			endTS:      "not-a-timestamp",
			wantError:  true,
			wantFields: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Diff(tt.startTS, tt.endTS)

			if (err != nil) != tt.wantError {
				t.Errorf("Diff() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				for _, field := range tt.wantFields {
					if _, ok := result[field]; !ok {
						t.Errorf("Diff() missing field %q in result", field)
					}
				}

				if tt.validateSign {
					secs := result["seconds"].(int64)
					isPositive := secs > 0
					if isPositive != tt.wantPositive {
						t.Errorf("Diff() seconds sign mismatch: got positive=%v, want positive=%v", isPositive, tt.wantPositive)
					}
				}
			}
		})
	}
}
