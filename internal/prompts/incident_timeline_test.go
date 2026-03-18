package prompts

import (
	"testing"
)

func TestIncidentTimeline(t *testing.T) {
	args := map[string]string{
		"events":           "10:00 - Service down, 10:15 - Alert triggered, 10:45 - Resolved",
		"source_timezones": "America/New_York,America/Chicago,America/Los_Angeles",
		"target_timezone":  "UTC",
	}

	result := IncidentTimeline(args)

	if result == "" {
		t.Error("IncidentTimeline() returned empty prompt")
	}

	// Verify key phrases are in the prompt
	expectedPhrases := []string{"incident", "timeline", "convert", "UTC", "time.diff"}
	for _, phrase := range expectedPhrases {
		if !contains(result, phrase) {
			t.Errorf("IncidentTimeline() missing expected phrase: %q", phrase)
		}
	}
}

func TestGetIncidentTimelineArgumentSchema(t *testing.T) {
	schema := IncidentTimelineArgumentSchema()

	if schema == nil {
		t.Error("GetIncidentTimelineArgumentSchema() returned nil")
		return
	}

	if schema["type"] != "object" {
		t.Errorf("GetIncidentTimelineArgumentSchema() type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Error("GetIncidentTimelineArgumentSchema() properties not a map")
		return
	}

	expectedProps := []string{"events", "source_timezones", "target_timezone"}
	for _, prop := range expectedProps {
		if _, ok := props[prop]; !ok {
			t.Errorf("GetIncidentTimelineArgumentSchema() missing property: %q", prop)
		}
	}
}
