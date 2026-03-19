package server

import (
	"log/slog"
	"time"
)

// LoggingMiddleware wraps a tool handler to add structured logging.
func LoggingMiddleware(method string, handler func(map[string]any) (map[string]any, error)) func(map[string]any) (map[string]any, error) {
	return func(args map[string]any) (map[string]any, error) {
		start := time.Now()
		result, err := handler(args)
		duration := time.Since(start)

		if err != nil {
			slog.Error(
				"Tool execution failed",
				"method", method,
				"duration_ms", duration.Milliseconds(),
				"error", err.Error(),
			)
		} else {
			slog.Info(
				"Tool execution succeeded",
				"method", method,
				"duration_ms", duration.Milliseconds(),
			)
		}

		return result, err
	}
}

// LoggingMiddlewarePrompt wraps a prompt handler to add structured logging.
func LoggingMiddlewarePrompt(method string, handler func(map[string]string) string) func(map[string]string) string {
	return func(args map[string]string) string {
		start := time.Now()
		result := handler(args)
		duration := time.Since(start)

		slog.Info(
			"Prompt generated",
			"method", method,
			"duration_ms", duration.Milliseconds(),
		)

		return result
	}
}
