//go:build integration

package integration

import (
	"testing"

	"github.com/tombenke/mcp-time-server/internal/server"
)

func TestToolErrorHandling(t *testing.T) {
	s := server.New()

	tests := []struct {
		name    string
		testFn  func() error
		wantErr bool
	}{
		{
			name: "Now with invalid timezone",
			testFn: func() error {
				_, err := s.ToolNow("Invalid/Zone", "")
				return err
			},
			wantErr: true,
		},
		{
			name: "Convert with invalid timestamp",
			testFn: func() error {
				_, err := s.ToolConvert("not-a-timestamp", "UTC", "UTC", "")
				return err
			},
			wantErr: true,
		},
		{
			name: "Diff with invalid start timestamp",
			testFn: func() error {
				_, err := s.ToolDiff("invalid", "2024-03-17T10:00:00Z")
				return err
			},
			wantErr: true,
		},
		{
			name: "TimezoneInfo with invalid timezone",
			testFn: func() error {
				_, err := s.ToolTimezoneInfo("Invalid/Zone", "")
				return err
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.testFn()
			if (err != nil) != tt.wantErr {
				t.Errorf("Error mismatch: got error %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
