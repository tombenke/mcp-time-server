package transport_test

import (
	"context"
	"testing"

	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	serverpkg "github.com/tombenke/mcp-time-server/internal/server"
	"github.com/tombenke/mcp-time-server/internal/transport"
)

func TestStreamableHTTPTransportSmoke(t *testing.T) {
	t.Parallel()

	s := serverpkg.New()
	httpServer := mcpserver.NewStreamableHTTPServer(
		s.MCP(),
		mcpserver.WithEndpointPath(transport.StreamableHTTPEndpointPath),
	)
	testServer := mcpserver.NewTestStreamableHTTPServer(s.MCP(), mcpserver.WithEndpointPath(transport.StreamableHTTPEndpointPath))
	defer testServer.Close()
	_ = httpServer

	client, err := client.NewStreamableHttpClient(testServer.URL + transport.StreamableHTTPEndpointPath)
	if err != nil {
		t.Fatalf("NewStreamableHttpClient() error = %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if _, err := client.Initialize(ctx, testInitializeRequest()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	result, err := client.CallTool(ctx, mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "time.convert",
			Arguments: map[string]any{
				"timestamp":     "2026-03-17T12:00:00Z",
				"from_timezone": "UTC",
				"to_timezone":   "Europe/Berlin",
			},
		},
	})
	if err != nil {
		t.Fatalf("CallTool() error = %v", err)
	}
	if result.IsError {
		t.Fatalf("CallTool() returned tool error: %+v", result.Content)
	}
}

func TestSSETransportSmoke(t *testing.T) {
	t.Parallel()

	s := serverpkg.New()
	testServer := mcpserver.NewTestServer(
		s.MCP(),
		mcpserver.WithStaticBasePath(transport.SSEBasePath),
	)
	defer testServer.Close()

	client, err := client.NewSSEMCPClient(testServer.URL + transport.SSEEndpointPath)
	if err != nil {
		t.Fatalf("NewSSEMCPClient() error = %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	if err := client.Start(ctx); err != nil {
		t.Fatalf("Start() error = %v", err)
	}

	if _, err := client.Initialize(ctx, testInitializeRequest()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	promptsResult, err := client.ListPrompts(ctx, mcp.ListPromptsRequest{})
	if err != nil {
		t.Fatalf("ListPrompts() error = %v", err)
	}
	if len(promptsResult.Prompts) != 3 {
		t.Fatalf("expected 3 prompts, got %d", len(promptsResult.Prompts))
	}
}

func testInitializeRequest() mcp.InitializeRequest {
	return mcp.InitializeRequest{
		Params: mcp.InitializeParams{
			ProtocolVersion: mcp.LATEST_PROTOCOL_VERSION,
			ClientInfo: mcp.Implementation{
				Name:    "transport-test-client",
				Version: "1.0.0",
			},
		},
	}
}
