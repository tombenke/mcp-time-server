---
title: MCP Time Server for Coding Agent Interoperability
version: 1.0
date_created: 2026-03-17
last_updated: 2026-03-17
owner: Agent Framework Team
tags: [app, design, infrastructure, mcp, time-server, interoperability]
---

# Introduction

This specification defines a production-grade MCP time server that is interoperable with modern coding agents and MCP clients, including Visual Studio Code with GitHub Copilot, Claude Code, Ollama-based clients, and the local Go agent framework in this repository. The server must implement real MCP protocol semantics and transport behavior, not a mock protocol subset. The solution includes tools, prompts, and supporting contracts required for reliable agent integration.

## 1. Purpose & Scope

Purpose:

- Define a standard-compliant MCP time server that exposes high-value time and timezone capabilities.
- Ensure protocol and transport compatibility across diverse MCP clients used in coding-agent workflows.
- Define tool and prompt contracts that are unambiguous for both humans and LLMs.

Scope:

- MCP protocol conformance for initialization, tool discovery, prompt discovery, and invocation.
- Required transport support for client interoperability.
- Tool set for time lookup, conversion, timezone metadata, and temporal calculations.
- Prompt set that complements tools for common planning and scheduling tasks.
- Testing, validation, and automation criteria for sustained compatibility.

Intended audience:

- Go developers implementing MCP servers.
- AI-agent platform engineers integrating external MCP servers.
- QA engineers validating MCP interoperability.

Assumptions:

- Server runtime is Go on macOS/Linux.
- Clients may use either legacy SSE or streamable HTTP MCP transports.
- Timezone data is sourced from IANA TZ database available to the host OS.

## 2. Definitions

| Term | Definition |
|------|------------|
| MCP | Model Context Protocol for tool, prompt, resource, and message interaction between clients and servers. |
| Stdio transport | MCP over process stdin/stdout JSON-RPC framing. |
| SSE transport | Legacy MCP transport where client opens SSE stream and posts messages to server endpoint. |
| Streamable HTTP | MCP transport using HTTP GET/POST/DELETE session semantics and optional SSE streaming. |
| IANA TZ | IANA Time Zone Database identifiers (example: `America/New_York`). |
| RFC3339 | Timestamp format used for machine-readable UTC and local times. |
| Prompt | MCP prompt definition returned by `prompts/list` and invoked by `prompts/get`. |
| Tool schema | JSON Schema contract describing tool input arguments. |

## 3. Requirements, Constraints & Guidelines

- **REQ-001**: The server SHALL implement MCP initialization flow (`initialize`, then `notifications/initialized`) according to current MCP spec behavior.
- **REQ-002**: The server SHALL implement `tools/list` and `tools/call` with valid JSON-RPC 2.0 envelopes.
- **REQ-003**: The server SHALL implement `prompts/list` and `prompts/get` so prompts are first-class capabilities beside tools.
- **REQ-004**: The server SHALL expose at least these tools: `time.now`, `time.convert`, `time.diff`, `time.timezone_info`, `time.list_timezones`.
- **REQ-005**: `time.now` SHALL support default timezone behavior and explicit timezone overrides.
- **REQ-006**: `time.convert` SHALL convert timestamps between valid IANA timezones with RFC3339 output.
- **REQ-007**: `time.diff` SHALL return signed and absolute differences between two timestamps in multiple units.
- **REQ-008**: `time.timezone_info` SHALL return UTC offset, DST state, and canonical timezone metadata for a target instant.
- **REQ-009**: `time.list_timezones` SHALL support region filtering and prefix search.
- **REQ-010**: Tool errors originating from business logic SHALL be returned as MCP tool errors (`isError=true`) with actionable text for LLM self-correction.
- **REQ-011**: The server SHALL support Stdio transport for local agent runtimes.
- **REQ-012**: The server SHALL support at least one HTTP-based MCP transport: Streamable HTTP and/or legacy SSE.
- **REQ-013**: If both Streamable HTTP and SSE are implemented, tool and prompt behavior SHALL be semantically identical across transports.
- **REQ-014**: SSE mode SHALL expose discoverable SSE and message endpoints and preserve session isolation.
- **REQ-015**: Streamable HTTP mode SHALL support session lifecycle and protocol-version headers required by compliant clients.
- **REQ-016**: Tool schemas SHALL be valid JSON Schema object contracts with explicit property types and required arrays.
- **REQ-017**: Tool and prompt names SHALL be stable and semantically versioned to avoid client regressions.
- **REQ-018**: Prompt definitions SHALL include deterministic prompt arguments and clear intended output format.
- **REQ-019**: The server SHALL include observability hooks for request method, duration, status, and error details.
- **REQ-020**: The server SHALL be compatible with the agent framework MCP host contract in this repository (SSE and stdio use cases).
- **REQ-021**: The server SHALL be fully compatible with the MCP Inspector tool and SHALL produce responses consumable by MCP Inspector without schema or protocol errors.
- **REQ-022**: Time calculations SHALL be deterministic for the same input and timezone database state.

- **SEC-001**: The server SHALL validate and normalize all user-supplied timezone identifiers.
- **SEC-002**: The server SHALL reject malformed timestamps with explicit validation messages.
- **SEC-003**: The server SHALL not execute arbitrary shell commands or file writes as part of tool behavior.
- **SEC-004**: Any network transport mode SHALL allow TLS termination strategy compatible with deployment requirements.

- **CON-001**: Implementation language SHALL be Go.
- **CON-002**: External dependencies SHALL be minimal and MCP focused.
- **CON-003**: Public tool and prompt contracts SHALL remain backward compatible across patch releases.

- **GUD-001**: Use package separation for transport wiring, tool logic, and prompt logic.
- **GUD-002**: Keep command entry points minimal and delegate behavior to internal packages.
- **GUD-003**: Prefer RFC3339 and ISO-8601-compatible fields for machine interoperability.

- **PAT-001**: Use Adapter pattern for transport abstraction (stdio, SSE, streamable HTTP).
- **PAT-002**: Use Contract-First pattern for tool/prompt JSON Schema design.
- **PAT-003**: Use Capability Advertisement pattern so clients can infer supported MCP features.

## 4. Interfaces & Data Contracts

### 4.1 Tool Definitions

Required tool contracts:

| Tool Name | Purpose | Required Arguments | Optional Arguments | Output Contract |
|-----------|---------|--------------------|--------------------|-----------------|
| `time.now` | Current time for a timezone | none | `timezone`, `format` | RFC3339 time, unix epoch, timezone id |
| `time.convert` | Convert timestamp across zones | `timestamp`, `from_timezone`, `to_timezone` | `format` | converted time, source, destination, offset delta |
| `time.diff` | Difference between two times | `start_timestamp`, `end_timestamp` | `unit` | signed and absolute delta in seconds/minutes/hours |
| `time.timezone_info` | Metadata for timezone at instant | `timezone` | `timestamp` | utc_offset, dst_active, abbreviation |
| `time.list_timezones` | Discover timezones | none | `region`, `prefix`, `limit` | array of matching timezone ids |

Example argument schema (`time.convert`):

```json
{
  "type": "object",
  "properties": {
    "timestamp": {"type": "string", "description": "Input time in RFC3339"},
    "from_timezone": {"type": "string", "description": "IANA timezone of input timestamp"},
    "to_timezone": {"type": "string", "description": "IANA timezone target"},
    "format": {"type": "string", "enum": ["rfc3339", "unix", "human"]}
  },
  "required": ["timestamp", "from_timezone", "to_timezone"]
}
```

### 4.2 Prompt Definitions

Required prompts beside tools:

| Prompt Name | Purpose | Arguments | Expected Prompt Output |
|-------------|---------|-----------|------------------------|
| `prompt.schedule_meeting` | Build cross-timezone meeting proposal | `participants`, `timezones`, `duration_minutes`, `date_window` | structured planning instructions with tool call suggestions |
| `prompt.incident_timeline` | Normalize event times in one canonical zone | `events`, `source_timezones`, `target_timezone` | ordered timeline transformation template |
| `prompt.timezone_decision` | Recommend timezone strategy for distributed teams | `team_regions`, `constraints` | trade-off analysis template |

### 4.3 MCP Capability Contract

The server initialize response SHALL advertise:

- tools capability (list and call)
- prompts capability (list and get)
- transport-specific capabilities where supported by implementation

### 4.4 Transport Contract

- Stdio mode: JSON-RPC over stdin/stdout, one process per server instance.
- HTTP mode: server SHALL implement Streamable HTTP and/or legacy SSE transport.
- SSE mode (if enabled): GET SSE stream endpoint and POST message endpoint with session correlation.
- Streamable HTTP mode (if enabled): MCP session creation and request handling over HTTP per protocol.

## 5. Acceptance Criteria

- **AC-001**: Given an MCP client initializes using valid protocol version, when `initialize` is sent, then the server returns capabilities including tools and prompts.
- **AC-002**: Given a client calls `tools/list`, when the request succeeds, then all required time tools are returned with valid JSON schemas.
- **AC-003**: Given a client calls `prompts/list`, when the request succeeds, then required prompts are returned with argument schemas.
- **AC-004**: Given an invalid timezone is supplied, when a time tool is invoked, then the response sets `isError=true` and includes a corrective message.
- **AC-005**: Given SSE transport is implemented and used, when a client starts stream and posts tool calls, then responses are delivered on the same logical session.
- **AC-006**: Given Streamable HTTP transport is implemented and used, when tool calls are issued, then responses follow streamable MCP semantics.
- **AC-007**: Given the server is checked with MCP Inspector, when validation runs, then no protocol or schema violations are reported.
- **AC-008**: Given the local agent framework connects via SSE or stdio, when discovering tools, then tool listing and invocation succeed end-to-end.

## 6. Test Automation Strategy

- **Test Levels**: Unit, Integration, End-to-End, Protocol Conformance.
- **Frameworks**: Go `testing`, `testify` (optional), MCP Inspector for protocol validation.
- **Test Data Management**: Deterministic fixed timestamps and timezone fixtures; no wall-clock assertions without tolerance windows.
- **CI/CD Integration**: GitHub Actions pipeline SHALL run `go test ./...`, protocol smoke tests, and optional inspector-based checks.
- **Coverage Requirements**: Core time conversion and validation logic SHALL maintain high coverage (target >= 85%).
- **Performance Testing**: Basic latency checks for tool invocation under concurrent sessions.

## 7. Rationale & Context

Coding agents rely on strict MCP behavior for reliable tool usage. A server that only partially implements MCP creates brittle integrations and increases hallucination risk during tool planning. Time servers are deceptively complex due to timezone rules, DST transitions, and transport diversity across clients. Including prompts beside tools enables better agent planning quality and reduces prompt-engineering duplication in client implementations.

## 8. Dependencies & External Integrations

### External Systems

- **EXT-001**: MCP Clients (VS Code GitHub Copilot, Claude Code, Ollama MCP clients, local agent framework) - protocol consumers.

### Third-Party Services

- **SVC-001**: None required for baseline operation; optional external timezone metadata sources are non-mandatory.

### Infrastructure Dependencies

- **INF-001**: Runtime environment capable of running a long-lived Go process and exposing stdio/HTTP endpoints.

### Data Dependencies

- **DAT-001**: IANA Time Zone Database - authoritative timezone and DST rules.

### Technology Platform Dependencies

- **PLT-001**: Go runtime with stable timezone handling and context-aware networking support.

### Compliance Dependencies

- **COM-001**: MCP protocol compliance requirements as defined by current and widely used client implementations.

## 9. Examples & Edge Cases

```code
// Edge case 1: DST spring-forward gap
// Input: 2026-03-08T02:30:00 America/New_York
// Expected: Validation error or normalized behavior documented by tool contract.

// Edge case 2: Ambiguous fall-back hour
// Input: 2026-11-01T01:30:00 America/New_York
// Expected: Result includes ambiguity handling strategy.

// Edge case 3: Invalid timezone
// Input timezone: "Mars/Olympus"
// Expected: isError=true with guidance to use valid IANA identifiers.

// Edge case 4: Leap second style input
// Input: 2016-12-31T23:59:60Z
// Expected: Explicit parse/validation result per parser constraints.
```

## 10. Validation Criteria

- All required tools and prompts are discoverable via MCP list methods.
- All tool input schemas are valid JSON Schema objects.
- Error responses for tool-level failures use MCP tool error semantics (`isError=true`).
- Stdio transport and at least one HTTP-based MCP transport (Streamable HTTP and/or SSE) pass interoperability smoke tests.
- Stdio transport works with local process clients.
- MCP Inspector validation reports no critical protocol or schema failures.
- End-to-end integration with local agent framework tool discovery and invocation is successful.

## 11. Related Specifications / Further Reading

- spec/spec-agent-framework.md
- <https://modelcontextprotocol.io>
- <https://github.com/mark3labs/mcp-go>
