package prompts

import (
	"testing"
)

func TestScheduleMeeting(t *testing.T) {
	args := map[string]string{
		"participants":     "Alice, Bob, Carol",
		"timezones":        "America/New_York,Europe/London,Asia/Tokyo",
		"duration_minutes": "90",
		"date_window":      "next 7 days",
	}

	result := ScheduleMeeting(args)

	if result == "" {
		t.Error("ScheduleMeeting() returned empty prompt")
	}

	// Verify key phrases are in the prompt
	expectedPhrases := []string{"meeting", "timezone", "participants", "time.now", "time.convert"}
	for _, phrase := range expectedPhrases {
		if !contains(result, phrase) {
			t.Errorf("ScheduleMeeting() missing expected phrase: %q", phrase)
		}
	}
}

func TestGetScheduleMeetingArgumentSchema(t *testing.T) {
	schema := ScheduleMeetingArgumentSchema()

	if schema == nil {
		t.Error("GetScheduleMeetingArgumentSchema() returned nil")
		return
	}

	if schema["type"] != "object" {
		t.Errorf("GetScheduleMeetingArgumentSchema() type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Error("GetScheduleMeetingArgumentSchema() properties not a map")
		return
	}

	expectedProps := []string{"participants", "timezones", "duration_minutes", "date_window"}
	for _, prop := range expectedProps {
		if _, ok := props[prop]; !ok {
			t.Errorf("GetScheduleMeetingArgumentSchema() missing property: %q", prop)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
