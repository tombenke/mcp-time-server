package prompts

import (
	"fmt"
)

// TimezoneDecision generates a prompt for timezone-related decision making.
func TimezoneDecision(args map[string]string) string {
	teamRegions := args["team_regions"]
	constraints := args["constraints"]

	if teamRegions == "" {
		teamRegions = "Global"
	}
	if constraints == "" {
		constraints = "Business hours priority"
	}

	prompt := fmt.Sprintf(`You are helping a distributed team make timezone-related decisions.

Team Regions: %s
Constraints: %s

To provide a comprehensive timezone analysis:
1. Use the 'time.now' tool to get current time in key regional hubs
2. Use the 'time.timezone_info' tool to understand UTC offsets and DST states
3. Use the 'time.list_timezones' tool with region filters
4. Use the 'time.convert' tool to simulate meetings at proposed times

Provide trade-off analysis based on constraints: %s
`, teamRegions, constraints, constraints)

	return prompt
}

// TimezoneDecisionArgumentSchema returns the argument schema for timezone_decision prompt.
func TimezoneDecisionArgumentSchema() map[string]any {
	return map[string]any{
		"type": "object",
		"properties": map[string]any{
			"team_regions": map[string]any{"type": "string", "description": "Geographic regions where team is distributed"},
			"constraints":  map[string]any{"type": "string", "description": "Business constraints or priorities"},
		},
		"required": []string{},
	}
}
