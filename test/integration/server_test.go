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
		t.Error("Capabilities doesnt
