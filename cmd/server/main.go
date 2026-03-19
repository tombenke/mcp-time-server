package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/tombenke/mcp-time-server/internal/server"
	"github.com/tombenke/mcp-time-server/internal/transport"
)

func main() {
	transportFlag := flag.String("transport", "stdio", "Transport type: stdio, sse, or http")
	addrFlag := flag.String("addr", "localhost:8080", "Server address for SSE and HTTP transports")
	flag.Parse()

	logWriter := io.Writer(os.Stderr)
	if *transportFlag == "stdio" {
		// Keep stdio transport strictly JSON-RPC and avoid bridge incompatibilities.
		logWriter = io.Discard
	}

	// Set up structured logging.
	slog.SetDefault(slog.New(
		slog.NewJSONHandler(logWriter, &slog.HandlerOptions{Level: slog.LevelInfo}),
	))

	slog.Info("Starting MCP Time Server", "transport", *transportFlag, "addr", *addrFlag)

	// Create the MCP server
	mcpServer := server.New()

	// Create a context that can be cancelled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-sigChan
		slog.Info("Received signal, shutting down", "signal", sig)
		cancel()
	}()

	var err error

	switch *transportFlag {
	case "stdio":
		err = transport.RunStdio(ctx, mcpServer)
	case "sse":
		err = transport.RunSSE(ctx, mcpServer, *addrFlag)
	case "http":
		err = transport.RunStreamableHTTP(ctx, mcpServer, *addrFlag)
	default:
		fmt.Fprintf(os.Stderr, "unsupported transport: %s\n", *transportFlag)
		os.Exit(1)
	}

	if err != nil && err != context.Canceled {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}

	slog.Info("Server stopped gracefully")
}
