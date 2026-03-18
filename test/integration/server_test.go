//go:build integration

package integration

import (
	"testing"

	"github.com/tombenke/mcp-time-server/internal/server"
)

func TestServerInitialization(t *testing.T) {
	s := server.New()

	if s == nil {
		t.Error("New() returned nil server")
		return
	}

	caps := s.GetCapabilities()

	if caps == nil {
		t.Error("GetCapabilities() returned nil")
		return
	}

	// Verify all tools are advertised
	tools, ok := caps["tools"].(map[string]any)
	if !ok {
		t.Error("capabilities do not include tools map")
		return
	}

	requiredTools := []string{
		"time.now",
		"time.convert",
		"time.diff",
		"time.timezone_info",
		"time.list_timezones",
	}

	for _, tool := range requiredTools {
		if _, exists := tools[tool]; !exists {
			t.Errorf("missing advertised tool: %s", tool)
		}
	}

	// Verify prompts are advertised
	prompts, ok := caps["prompts"].(map[string]any)
	if !ok {
		t.Error("capabilities do not include prompts map")
		return
	}

	requiredPrompts := []string{
		"prompt.schedule_meeting",
		"prompt.incident_timeline",
		"prompt.timezone_decision",
	}

	for _, prompt := range requiredPrompts {
		if _, exists := prompts[prompt]; !exists {
			t.Errorf("missing advertised prompt: %s", prompt)
		}
	}
}
