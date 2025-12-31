package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/config"
	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/mcp"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Convert string log level to slog.Level
	var level slog.Level
	switch strings.ToLower(cfg.LogLevel) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	// Setup logger with level
	// Use stderr for logging so stdout is available for stdio transport
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)
	slog.SetDefault(logger)

	// Mask tokens for logging
	maskToken := func(token string) string {
		if len(token) <= 4 {
			return "****"
		}
		return token[:2] + strings.Repeat("*", len(token)-4) + token[len(token)-2:]
	}

	slog.Info("Starting Paperless MCP Server",
		"paperless_url", cfg.PaperlessURL,
		"mcp_transport", cfg.MCPTransport,
		"log_level", cfg.LogLevel,
		"paperless_token", maskToken(cfg.PaperlessToken),
		"mcp_auth_token", maskToken(cfg.MCPAuthToken),
		"mcp_http_port", cfg.MCPHTTPPort,
	)

	// Create MCP server
	mcpServer, err := mcp.New(cfg)
	if err != nil {
		slog.Error("Failed to create MCP server", "error", err)
		os.Exit(1)
	}

	// Create context with cancellation for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		slog.Info("Shutdown signal received", "signal", sig)
		cancel()
	}()

	// Start server with appropriate transport
	var serverErr error
	switch cfg.MCPTransport {
	case "stdio":
		slog.Info("Starting with stdio transport")
		serverErr = mcpServer.StartStdio(ctx)
	case "http":
		slog.Info("Starting with HTTP transport", "port", cfg.MCPHTTPPort)
		serverErr = mcpServer.StartHTTP(ctx)
	default:
		slog.Error("Invalid transport mode", "transport", cfg.MCPTransport)
		os.Exit(1)
	}

	if serverErr != nil {
		slog.Error("Server error", "error", serverErr)
		os.Exit(1)
	}

	slog.Info("Server shutdown complete")
}
