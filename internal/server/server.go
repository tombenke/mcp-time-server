package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	mcpserver "github.com/mark3labs/mcp-go/server"
	"github.com/tombenke/mcp-time-server/internal/prompts"
	"github.com/tombenke/mcp-time-server/internal/tools"
)

const serverInstructions = "Use the time.* tools for timezone-aware calculations and the prompt.* templates for scheduling, incident timeline normalization, and timezone decision support."

// Server manages the MCP server instance and tool/prompt registrations.
type Server struct {
	name    string
	version string
	mcp     *mcpserver.MCPServer
}

// New creates a new MCP server with all tools and prompts registered.
func New() *Server {
	s := &Server{
		name:    "mcp-time-server",
		version: "1.0.0",
	}
	s.mcp = s.newMCPServer()

	slog.Info("MCP server initialized", "name", s.name, "version", s.version)

	return s
}

// MCP returns the registered MCP server instance.
func (s *Server) MCP() *mcpserver.MCPServer {
	return s.mcp
}

// ToolNow handles the time.now tool request.
func (s *Server) ToolNow(timezone, format string) (map[string]any, error) {
	start := time.Now()
	result, err := tools.Now(timezone, format)
	duration := time.Since(start)

	if err != nil {
		slog.Error("Tool execution failed", "method", "time.now", "duration_ms", duration.Milliseconds(), "error", err)
	} else {
		slog.Info("Tool execution succeeded", "method", "time.now", "duration_ms", duration.Milliseconds())
	}

	return result, err
}

// ToolConvert handles the time.convert tool request.
func (s *Server) ToolConvert(timestamp, fromTZ, toTZ, format string) (map[string]any, error) {
	start := time.Now()
	result, err := tools.Convert(timestamp, fromTZ, toTZ, format)
	duration := time.Since(start)

	if err != nil {
		slog.Error("Tool execution failed", "method", "time.convert", "duration_ms", duration.Milliseconds(), "error", err)
	} else {
		slog.Info("Tool execution succeeded", "method", "time.convert", "duration_ms", duration.Milliseconds())
	}

	return result, err
}

// ToolDiff handles the time.diff tool request.
func (s *Server) ToolDiff(start, end string) (map[string]any, error) {
	startTime := time.Now()
	result, err := tools.Diff(start, end)
	duration := time.Since(startTime)

	if err != nil {
		slog.Error("Tool execution failed", "method", "time.diff", "duration_ms", duration.Milliseconds(), "error", err)
	} else {
		slog.Info("Tool execution succeeded", "method", "time.diff", "duration_ms", duration.Milliseconds())
	}

	return result, err
}

// ToolTimezoneInfo handles the time.timezone_info tool request.
func (s *Server) ToolTimezoneInfo(timezone, timestamp string) (map[string]any, error) {
	start := time.Now()
	result, err := tools.TimezoneInfo(timezone, timestamp)
	duration := time.Since(start)

	if err != nil {
		slog.Error("Tool execution failed", "method", "time.timezone_info", "duration_ms", duration.Milliseconds(), "error", err)
	} else {
		slog.Info("Tool execution succeeded", "method", "time.timezone_info", "duration_ms", duration.Milliseconds())
	}

	return result, err
}

// ToolListTimezones handles the time.list_timezones tool request.
func (s *Server) ToolListTimezones(region, prefix string, limit int) ([]string, error) {
	start := time.Now()
	result, err := tools.ListTimezones(region, prefix, limit)
	duration := time.Since(start)

	if err != nil {
		slog.Error("Tool execution failed", "method", "time.list_timezones", "duration_ms", duration.Milliseconds(), "error", err)
	} else {
		slog.Info("Tool execution succeeded", "method", "time.list_timezones", "duration_ms", duration.Milliseconds())
	}

	return result, err
}

// PromptScheduleMeeting handles the prompt.schedule_meeting prompt request.
func (s *Server) PromptScheduleMeeting(args map[string]string) string {
	start := time.Now()
	result := prompts.ScheduleMeeting(args)
	duration := time.Since(start)

	slog.Info("Prompt generated", "method", "prompt.schedule_meeting", "duration_ms", duration.Milliseconds())
	return result
}

// PromptIncidentTimeline handles the prompt.incident_timeline prompt request.
func (s *Server) PromptIncidentTimeline(args map[string]string) string {
	start := time.Now()
	result := prompts.IncidentTimeline(args)
	duration := time.Since(start)

	slog.Info("Prompt generated", "method", "prompt.incident_timeline", "duration_ms", duration.Milliseconds())
	return result
}

// PromptTimezoneDecision handles the prompt.timezone_decision prompt request.
func (s *Server) PromptTimezoneDecision(args map[string]string) string {
	start := time.Now()
	result := prompts.TimezoneDecision(args)
	duration := time.Since(start)

	slog.Info("Prompt generated", "method", "prompt.timezone_decision", "duration_ms", duration.Milliseconds())
	return result
}

// GetCapabilities returns the server capabilities.
func (s *Server) GetCapabilities() map[string]any {
	return map[string]any{
		"tools": map[string]any{
			"time.now":            "Current time in timezone",
			"time.convert":        "Convert timestamp across timezones",
			"time.diff":           "Calculate time difference",
			"time.timezone_info":  "Get timezone information",
			"time.list_timezones": "List available timezones",
		},
		"prompts": map[string]any{
			"prompt.schedule_meeting":  "Schedule meeting across timezones",
			"prompt.incident_timeline": "Normalize incident timeline",
			"prompt.timezone_decision": "Analyze timezone decisions",
		},
	}
}

func (s *Server) newMCPServer() *mcpserver.MCPServer {
	mcpSrv := mcpserver.NewMCPServer(
		s.name,
		s.version,
		mcpserver.WithToolCapabilities(true),
		mcpserver.WithPromptCapabilities(true),
		mcpserver.WithInstructions(serverInstructions),
	)

	s.registerTools(mcpSrv)
	s.registerPrompts(mcpSrv)

	return mcpSrv
}

func (s *Server) registerTools(mcpSrv *mcpserver.MCPServer) {
	mcpSrv.AddTool(
		mcp.NewTool(
			"time.now",
			mcp.WithDescription("Current time in a timezone"),
			mcp.WithString("timezone", mcp.Description("IANA timezone identifier, defaults to UTC")),
			mcp.WithString("format", mcp.Description("Optional output format hint, currently informational only")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_ = ctx
			result, err := s.ToolNow(request.GetString("timezone", "UTC"), request.GetString("format", ""))
			return toolResult(result, err), nil
		},
	)

	mcpSrv.AddTool(
		mcp.NewTool(
			"time.convert",
			mcp.WithDescription("Convert an RFC3339 timestamp between timezones"),
			mcp.WithString("timestamp", mcp.Required(), mcp.Description("RFC3339 timestamp to convert")),
			mcp.WithString("from_timezone", mcp.Description("Source IANA timezone identifier, defaults to UTC")),
			mcp.WithString("to_timezone", mcp.Description("Destination IANA timezone identifier, defaults to UTC")),
			mcp.WithString("format", mcp.Description("Optional output format hint, currently informational only")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_ = ctx
			result, err := s.ToolConvert(
				request.GetString("timestamp", ""),
				request.GetString("from_timezone", "UTC"),
				request.GetString("to_timezone", "UTC"),
				request.GetString("format", ""),
			)
			return toolResult(result, err), nil
		},
	)

	mcpSrv.AddTool(
		mcp.NewTool(
			"time.diff",
			mcp.WithDescription("Calculate the signed and absolute difference between two RFC3339 timestamps"),
			mcp.WithString("start_timestamp", mcp.Required(), mcp.Description("Start RFC3339 timestamp")),
			mcp.WithString("end_timestamp", mcp.Required(), mcp.Description("End RFC3339 timestamp")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_ = ctx
			result, err := s.ToolDiff(
				request.GetString("start_timestamp", ""),
				request.GetString("end_timestamp", ""),
			)
			return toolResult(result, err), nil
		},
	)

	mcpSrv.AddTool(
		mcp.NewTool(
			"time.timezone_info",
			mcp.WithDescription("Return UTC offset, abbreviation, and DST state for a timezone"),
			mcp.WithString("timezone", mcp.Description("IANA timezone identifier, defaults to UTC")),
			mcp.WithString("timestamp", mcp.Description("Optional RFC3339 timestamp to evaluate instead of now")),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_ = ctx
			result, err := s.ToolTimezoneInfo(
				request.GetString("timezone", "UTC"),
				request.GetString("timestamp", ""),
			)
			return toolResult(result, err), nil
		},
	)

	mcpSrv.AddTool(
		mcp.NewTool(
			"time.list_timezones",
			mcp.WithDescription("List available IANA timezones with optional region and prefix filtering"),
			mcp.WithString("region", mcp.Description("Optional region prefix such as Europe or America")),
			mcp.WithString("prefix", mcp.Description("Optional case-insensitive substring filter")),
			mcp.WithNumber("limit", mcp.Description("Maximum number of timezones to return"), mcp.DefaultNumber(100), mcp.Min(1), mcp.Max(500)),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			_ = ctx
			result, err := s.ToolListTimezones(
				request.GetString("region", ""),
				request.GetString("prefix", ""),
				request.GetInt("limit", 100),
			)
			if err != nil {
				return toolErrorResult(err), nil
			}
			return mcp.NewToolResultStructuredOnly(map[string]any{"timezones": result}), nil
		},
	)
}

func (s *Server) registerPrompts(mcpSrv *mcpserver.MCPServer) {
	mcpSrv.AddPrompt(
		mcp.NewPrompt(
			"prompt.schedule_meeting",
			mcp.WithPromptDescription("Generate a meeting-planning prompt for distributed teams"),
			mcp.WithArgument("participants", mcp.ArgumentDescription("Team member names or count")),
			mcp.WithArgument("timezones", mcp.ArgumentDescription("Comma-separated IANA timezones"), mcp.RequiredArgument()),
			mcp.WithArgument("duration_minutes", mcp.ArgumentDescription("Meeting duration in minutes")),
			mcp.WithArgument("date_window", mcp.ArgumentDescription("Preferred date range, for example next 7 days")),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			_ = ctx
			return newPromptResult(
				"Schedule meeting across timezones",
				s.PromptScheduleMeeting(request.Params.Arguments),
			), nil
		},
	)

	mcpSrv.AddPrompt(
		mcp.NewPrompt(
			"prompt.incident_timeline",
			mcp.WithPromptDescription("Generate a prompt for normalizing incident timelines across timezones"),
			mcp.WithArgument("events", mcp.ArgumentDescription("Incident events with timestamps and descriptions"), mcp.RequiredArgument()),
			mcp.WithArgument("source_timezones", mcp.ArgumentDescription("Comma-separated source IANA timezones"), mcp.RequiredArgument()),
			mcp.WithArgument("target_timezone", mcp.ArgumentDescription("Target IANA timezone for the normalized timeline")),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			_ = ctx
			return newPromptResult(
				"Normalize incident timeline",
				s.PromptIncidentTimeline(request.Params.Arguments),
			), nil
		},
	)

	mcpSrv.AddPrompt(
		mcp.NewPrompt(
			"prompt.timezone_decision",
			mcp.WithPromptDescription("Generate a timezone decision-analysis prompt"),
			mcp.WithArgument("team_regions", mcp.ArgumentDescription("Geographic regions where the team is distributed")),
			mcp.WithArgument("constraints", mcp.ArgumentDescription("Business constraints or priorities")),
		),
		func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
			_ = ctx
			return newPromptResult(
				"Analyze timezone decisions",
				s.PromptTimezoneDecision(request.Params.Arguments),
			), nil
		},
	)
}

func newPromptResult(description, prompt string) *mcp.GetPromptResult {
	return &mcp.GetPromptResult{
		Description: description,
		Messages: []mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.RoleUser, mcp.NewTextContent(prompt)),
		},
	}
}

func toolResult(result any, err error) *mcp.CallToolResult {
	if err != nil {
		return toolErrorResult(err)
	}

	return mcp.NewToolResultStructuredOnly(result)
}

func toolErrorResult(err error) *mcp.CallToolResult {
	result := mcp.NewToolResultText(err.Error())
	result.IsError = true
	return result
}
