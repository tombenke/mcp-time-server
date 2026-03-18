package prompts

import (
	"testing"
)

func TestTimezoneDecision(t *testing.T) {
	args := map[string]string{
		"team_regions": "North America, Europe, Asia",
		"constraints":  "Maximize working hours overlap, minimize late nights",
	}

	result := TimezoneDecision(args)

	if result == "" {
		t.Error("TimezoneDecision() returned empty prompt")
	}

	// Verify key phrases are in the prompt
	expectedPhrases := []string{"timezone", "decision", "team", "time.now", "time.timezone_info"}
	for _, phrase := range expectedPhrases {
		if !contains(result, phrase) {
			t.Errorf("TimezoneDecision() missing expected phrase: %q", phrase)
		}
	}
}

func TestGetTimezoneDecisionArgumentSchema(t *testing.T) {
	schema := TimezoneDecisionArgumentSchema()

	if schema == nil {
		t.Error("GetTimezoneDecisionArgumentSchema() returned nil")
		return
	}

	if schema["type"] != "object" {
		t.Errorf("GetTimezoneDecisionArgumentSchema() type = %v, want object", schema["type"])
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		t.Error("GetTimezoneDecisionArgumentSchema() properties not a map")
		return
	}

	expectedProps := []string{"team_regions", "constraints"}
	for _, prop := range expectedProps {
		if _, ok := props[prop]; !ok {
			t.Errorf("GetTimezoneDecisionArgumentSchema() missing property: %q", prop)
		}
	}
}
