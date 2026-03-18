package prompts

import (
	"fmt"
)

// IncidentTimeline generates a prompt for normalizing an incident timeline across timezones.
func IncidentTimeline(args map[string]string) string {
	events := args["events"]
	sourceTimezones := args["source_timezones"]
	targetTimezone := args["target_timezone"]

	if targetTimezone == "" {
		targetTimezone = "UTC"
	}

	prompt := fmt.Sprintf(`You are analyzing an incident that occurred across multiple timezones and need to build a normalized timeline.

Events:
%s

Source timezones:
%s

Target timezone:
%s

To create an accurate incident timeline:
1. Use the 'time.convert' tool to convert each event timestamp to the target timezone
2. Use the 'time.diff' tool to calculate time deltas between events
3. Use the 'time.timezone_info' tool to check for DST or timezone offset impacts
`, events, sourceTimezones, targetTimezone)

	return prompt
}

// IncidentTimelineArgumentSchema returns the argument schema for incident_timeline prompt.
func IncidentTimelineArgumentSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"events":           map[string]any{"type": "string", "description": "List of incident events with timestamps and descriptions"},
			"source_timezones": map[string]any{"type": "string", "description": "Comma-separated IANA timezones for source events"},
			"target_timezone":  map[string]any{"type": "string", "description": "Target IANA timezone for normalized timeline (default: UTC)"},
		},
		"required": []string{"events", "source_timezones"},
	}
}
