# MCP Time Server

[![Actions Status](https://github.com/tombenke/mcp-time-server/workflows/ci/badge.svg)](https://github.com/tombenke/mcp-time-server)

Go-based time and timezone server intended for MCP-style tool and prompt workflows.

## Features

- Time tools
  - `time.now`: current time in a timezone
  - `time.convert`: timestamp conversion between timezones
  - `time.diff`: signed and absolute time difference
  - `time.timezone_info`: UTC offset, abbreviation, DST indicator
  - `time.list_timezones`: filtered timezone listing with limit
- Prompt templates
  - `prompt.schedule_meeting`
  - `prompt.incident_timeline`
  - `prompt.timezone_decision`
- Transport entry points
  - `stdio`
  - `sse`
  - `http`
- Structured logging via `log/slog`
- Unit and integration-style tests in place for current implementation

## Project Layout

- [cmd/server/main.go](cmd/server/main.go): application entry point and transport selection
- [internal/tools](internal/tools): pure time business logic
- [internal/prompts](internal/prompts): prompt generation and schemas
- [internal/server](internal/server): server facade and capability advertisement
- [internal/transport/transport.go](internal/transport/transport.go): transport runner functions
- [test/integration](test/integration): integration-oriented tests

## Requirements

- Go 1.21+ (project currently uses Go 1.25 in [go.mod](go.mod))

## Build

```bash
go mod tidy
go build ./...
```

## Run

Start with default transport (`stdio`):

```bash
go run ./cmd/server
```

Run with explicit transport:

```bash
go run ./cmd/server --transport stdio
go run ./cmd/server --transport sse --addr :8080
go run ./cmd/server --transport http --addr :8080
```

Build a standalone binary (recommended for Inspector stdio mode):

```bash
go build -o ./server ./cmd/server
./server --transport stdio
```

Flags:

- `--transport`: `stdio`, `sse`, `http` (default: `stdio`)
- `--addr`: bind address for network transports (default: `localhost:8080`)

Stop the process with `Ctrl+C` for graceful shutdown.

## Using with mcpinspector

### Stdio Transport (Recommended for Quick Testing)

1. In mcpinspector, select **Stdio** mode
2. Configure:
   - **Command**: `./server` (preferred) or `go`
   - **Args**:
     - for `./server`: `--transport stdio`
     - for `go`: `run ./cmd/server --transport stdio`
   - **Working Directory**: repository root
3. Click connect

Stdio troubleshooting sample:

```bash
go build -o ./server ./cmd/server
printf '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"probe","version":"0.1"}}}\n' | ./server --transport stdio
```

If stdio is configured correctly, the command above returns a JSON-RPC `initialize` result.

### SSE Transport (Browser/Web-based)

1. Start the server:

   ```bash
   go run ./cmd/server --transport sse --addr :8080
   ```

2. In mcpinspector, select **SSE** mode
3. Configure:
   - **SSE URI**: `http://127.0.0.1:8080/mcp/sse`
   - Or use `http://localhost:8080/mcp/sse` if connecting locally
4. Click connect

**Note:** The server advertises relative message endpoints, so browser origin mismatch errors are avoided.

SSE probe sample:

```bash
curl -i --max-time 3 http://127.0.0.1:8080/mcp/sse
```

### Streamable HTTP Transport

1. Start the server:

   ```bash
   go run ./cmd/server --transport http --addr :8080
   ```

2. In mcpinspector, select **Streamable HTTP** mode
3. Configure:
   - **URI**: `http://127.0.0.1:8080/mcp`
   - Or `http://localhost:8080/mcp`
4. Click connect

Streamable HTTP initialize sample:

```bash
curl -i -X POST 'http://127.0.0.1:8080/mcp' \
   -H 'Content-Type: application/json' \
   -H 'MCP-Protocol-Version: 2025-03-26' \
   --data '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-03-26","capabilities":{},"clientInfo":{"name":"curl","version":"0.1"}}}'
```

### Connection Details

- **CORS**: Both SSE (`/mcp/sse`, `/mcp/message`) and streamable HTTP (`/mcp`) include browser-compatible CORS and preflight handling
- **Stdio Output Safety**: In stdio mode, application logs are disabled so stdout remains clean MCP JSON-RPC output
- **Session Management**: Automatically handled by the MCP protocol layer
- **Protocol Version**: Negotiated during initialization by mcpinspector

## Test

```bash
go test ./...
```

## Current Behavior Notes

- Tool logic in [internal/tools](internal/tools) is functional and tested.
- Prompt generators in [internal/prompts](internal/prompts) are functional and tested.
- The server is wired to `github.com/mark3labs/mcp-go` and exposes working `stdio`, streamable HTTP, and SSE transports.
- Streamable HTTP uses the `/mcp` endpoint and SSE uses `/mcp/sse` plus `/mcp/message`.

## Next Steps

- Split transport runners into dedicated files for stdio, SSE, and streamable HTTP.
- Add protocol-level SSE/HTTP integration tests.
