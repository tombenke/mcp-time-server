package transport

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	mcpserver "github.com/mark3labs/mcp-go/server"
	appserver "github.com/tombenke/mcp-time-server/internal/server"
)

const (
	StreamableHTTPEndpointPath = "/mcp"
	SSEBasePath                = "/mcp"
	SSEEndpointPath            = "/mcp/sse"
	SSEMessageEndpointPath     = "/mcp/message"
	shutdownTimeout            = 5 * time.Second
)

// RunStdio starts the MCP server using Stdio transport.
func RunStdio(ctx context.Context, s *appserver.Server) error {
	slog.Info("Starting MCP server on stdio transport")
	return mcpserver.NewStdioServer(s.MCP()).Listen(ctx, os.Stdin, os.Stdout)
}

// RunSSE starts the MCP server using Server-Sent Events (SSE) transport.
func RunSSE(ctx context.Context, s *appserver.Server, addr string) error {
	listenAddr := normalizeListenAddr(addr)

	httpServer := &http.Server{}
	sseServer := mcpserver.NewSSEServer(
		s.MCP(),
		mcpserver.WithStaticBasePath(SSEBasePath),
		mcpserver.WithHTTPServer(httpServer),
	)

	mux := http.NewServeMux()
	mux.HandleFunc(SSEEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		sseServer.SSEHandler().ServeHTTP(w, r)
	})

	mux.HandleFunc(SSEMessageEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		sseServer.MessageHandler().ServeHTTP(w, r)
	})
	httpServer.Handler = mux

	slog.Info(
		"Starting MCP server on SSE transport",
		"addr", listenAddr,
		"sse_endpoint", SSEEndpointPath,
		"message_endpoint", SSEMessageEndpointPath,
	)

	return runHTTPServer(ctx, func() error {
		return sseServer.Start(listenAddr)
	}, func(shutdownCtx context.Context) error {
		return sseServer.Shutdown(shutdownCtx)
	})
}

// RunStreamableHTTP starts the MCP server using Streamable HTTP transport.
func RunStreamableHTTP(ctx context.Context, s *appserver.Server, addr string) error {
	listenAddr := normalizeListenAddr(addr)
	httpServer := &http.Server{}
	server := mcpserver.NewStreamableHTTPServer(
		s.MCP(),
		mcpserver.WithEndpointPath(StreamableHTTPEndpointPath),
		mcpserver.WithStreamableHTTPServer(httpServer),
	)

	mux := http.NewServeMux()
	mux.HandleFunc(StreamableHTTPEndpointPath, func(w http.ResponseWriter, r *http.Request) {
		addCORSHeaders(w, r)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		server.ServeHTTP(w, r)
	})
	httpServer.Handler = mux

	slog.Info(
		"Starting MCP server on streamable HTTP transport",
		"addr", listenAddr,
		"endpoint", StreamableHTTPEndpointPath,
	)

	return runHTTPServer(ctx, func() error {
		return server.Start(listenAddr)
	}, func(shutdownCtx context.Context) error {
		return server.Shutdown(shutdownCtx)
	})
}

func runHTTPServer(ctx context.Context, start func() error, shutdown func(context.Context) error) error {
	shutdownDone := make(chan struct{})

	go func() {
		defer close(shutdownDone)
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("Transport shutdown failed", "error", err)
		}
	}()

	err := start()
	<-shutdownDone

	if err == nil {
		return nil
	}
	if errors.Is(err, http.ErrServerClosed) && ctx.Err() != nil {
		return ctx.Err()
	}

	return err
}

func addCORSHeaders(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", requestedOrDefaultHeaders(r))
	w.Header().Set("Access-Control-Expose-Headers", "Mcp-Session-Id, MCP-Session-Id")
}

func requestedOrDefaultHeaders(r *http.Request) string {
	requested := strings.TrimSpace(r.Header.Get("Access-Control-Request-Headers"))
	if requested != "" {
		return requested
	}

	return "Content-Type, Accept, Last-Event-ID, Mcp-Session-Id, MCP-Session-Id, MCP-Protocol-Version"
}

func normalizeListenAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if addr == "" {
		return "localhost:8080"
	}
	if strings.HasPrefix(addr, "http://") {
		return strings.TrimPrefix(addr, "http://")
	}
	if strings.HasPrefix(addr, "https://") {
		return strings.TrimPrefix(addr, "https://")
	}

	return addr
}

func baseURLFromAddr(addr string) string {
	addr = strings.TrimSpace(addr)
	if strings.HasPrefix(addr, "http://") || strings.HasPrefix(addr, "https://") {
		u, err := url.Parse(strings.TrimSuffix(addr, "/"))
		if err == nil && u.Host != "" {
			host, port, splitErr := net.SplitHostPort(u.Host)
			if splitErr == nil {
				host = normalizeAdvertisedHost(host)
				u.Host = net.JoinHostPort(host, port)
			} else {
				u.Host = normalizeAdvertisedHost(u.Host)
			}
			return u.String()
		}

		return strings.TrimSuffix(addr, "/")
	}

	listenAddr := normalizeListenAddr(addr)
	if strings.HasPrefix(listenAddr, ":") {
		return "http://localhost" + listenAddr
	}

	host, port, err := net.SplitHostPort(listenAddr)
	if err == nil {
		host = normalizeAdvertisedHost(host)
		return "http://" + net.JoinHostPort(host, port)
	}

	host = normalizeAdvertisedHost(listenAddr)
	if host == "" {
		host = "localhost"
	}

	return "http://" + host
}

func normalizeAdvertisedHost(host string) string {
	host = strings.TrimSpace(host)
	host = strings.Trim(host, "[]")

	if host == "" || host == "0.0.0.0" || host == "::" {
		return "localhost"
	}

	return host
}
