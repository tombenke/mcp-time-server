package prompts

import (
	"fmt"
)

// ScheduleMeetingArgs defines the arguments for the schedule_meeting prompt.
type ScheduleMeetingArgs struct {
	Participants    string
	Timezones       string
	DurationMinutes string
	DateWindow      string
}

// ScheduleMeeting generates a prompt for scheduling a meeting across timezones.
func ScheduleMeeting(args map[string]string) string {
	participants := args["participants"]
	timezones := args["timezones"]
	duration := args["duration_minutes"]
	dateWindow := args["date_window"]

	if participants == "" {
		participants = "Team members"
	}
	if duration == "" {
		duration = "60"
	}
	if dateWindow == "" {
		dateWindow = "next 7 days"
	}

	prompt := fmt.Sprintf(`You are scheduling a meeting for participants: %s
Timezones: %s
Duration: %s minutes
Date Window: %s

To find the optimal meeting time:
1. Use 'time.now' tool for each timezone
2. Use 'time.convert' tool to convert proposed times
3. Use 'time.timezone_info' tool to check for DST impacts

Provide 3 alternative time slots.
`, participants, timezones, duration, dateWindow)

	return prompt
}

// ScheduleMeetingArgumentSchema returns the argument schema for schedule_meeting prompt.
func ScheduleMeetingArgumentSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"participants":     map[string]any{"type": "string", "description": "Team member names or count"},
			"timezones":        map[string]any{"type": "string", "description": "Comma-separated list of IANA timezones"},
			"duration_minutes": map[string]any{"type": "string", "description": "Meeting duration in minutes"},
			"date_window":      map[string]any{"type": "string", "description": "Preferred date range (e.g., next 7 days)"},
		},
		"required": []string{"timezones"},
	}

}
