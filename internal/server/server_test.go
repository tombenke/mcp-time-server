package server_test

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	serverpkg "github.com/tombenke/mcp-time-server/internal/server"
)

func TestNewRegistersToolsAndPrompts(t *testing.T) {
	t.Parallel()

	s := serverpkg.New()
	client, err := client.NewInProcessClient(s.MCP())
	if err != nil {
		t.Fatalf("NewInProcessClient() error = %v", err)
	}
	defer func() { _ = client.Close() }()

	ctx := context.Background()
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if _, err := client.Initialize(ctx, testInitializeRequest()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	toolsResult, err := client.ListTools(ctx, mcp.ListToolsRequest{})
	if err != nil {
		t.Fatalf("ListTools() error = %v", err)
	}
	if len(toolsResult.Tools) != 5 {
		t.Fatalf("expected 5 tools, got %d", len(toolsResult.Tools))
	}

	promptsResult, err := client.ListPrompts(ctx, mcp.ListPromptsRequest{})
	if err != nil {
		t.Fatalf("ListPrompts() error = %v", err)
	}
	if len(promptsResult.Prompts) != 3 {
		t.Fatalf("expected 3 prompts, got %d", len(promptsResult.Prompts))
	}

	toolCall, err := client.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "time.now",
			Arguments: map[string]any{"timezone": "UTC"},
		},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}
	if toolCall.IsError {
		t.Fatalf("CallTool() returned tool error: %+v", toolCall.Content)
	}

	promptResult, err := client.GetPrompt(ctx, mcp.GetPromptRequest{
		Params: mcp.GetPromptParams{
			Name: "prompt.schedule_meeting",
			Arguments: map[string]string{
				"timezones": "UTC,Europe/Berlin",
			},
		},
	})
	if err != nil {
		t.Fatalf("GetPrompt() error = %v", err)
	}
	if len(promptResult.Messages) != 1 {
		t.Fatalf("expected 1 prompt message, got %d", len(promptResult.Messages))
	}
}

func testInitializeRequest() mcp.InitializeRequest {
	return mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "server-test-client",
				Version: "1.0.0",
			},
		},
	}
}
