# AGENT.md

This repository contains an MCP time server written in Go.

## Project Overview

- Entry point: cmd/server/main.go
- MCP server wiring: internal/server/server.go
- Transport adapters: internal/transport/transport.go
- Tool logic: internal/tools/
- Prompt logic: internal/prompts/
- Integration tests: test/integration/

## Common Commands

- Build:
  - go build ./...
- Test:
  - go test ./...
- Run stdio transport:
  - go run ./cmd/server --transport stdio
- Run SSE transport:
  - go run ./cmd/server --transport sse --addr :8080
- Run streamable HTTP transport:
  - go run ./cmd/server --transport http --addr :8080

## Transport Notes

- Stdio mode must keep stdout protocol-clean JSON-RPC.
- SSE endpoint path: /mcp/sse
- SSE message endpoint path: /mcp/message
- Streamable HTTP endpoint path: /mcp
- Browser clients require CORS/preflight support.

## Development Guidelines

- Keep cmd/server/main.go minimal.
- Put business logic in internal/tools and internal/prompts.
- Preserve public MCP tool and prompt names for compatibility.
- Add tests for new tools, prompts, and transport behavior changes.

## Validation Checklist

Before opening a PR:

1. Run gofmt on changed Go files.
2. Run go build ./...
3. Run go test ./...
4. Verify transport changes against at least one MCP client.
