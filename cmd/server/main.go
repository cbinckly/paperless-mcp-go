package main

import (
    "fmt"
    "log/slog"
    "os"
    "strings"

    "git.binckly.ca/cbinckly/paperless-mcp-go/internal/config"
)

func main() {
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
    handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
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

    // TODO: Server initialization will be added in Task 3
}
