package config

import (
    "errors"
    "fmt"
    "os"
    "strings"
)

// Environment variable name constants
const (
    EnvPaperlessURL    = "PAPERLESS_URL"
    EnvPaperlessToken  = "PAPERLESS_TOKEN"
    EnvMCPAuthToken    = "MCP_AUTH_TOKEN"
    EnvLogLevel        = "LOG_LEVEL"
    EnvMCPTransport    = "MCP_TRANSPORT"
    EnvMCPHTTPPort     = "MCP_HTTP_PORT"
)

// Default values
const (
    DefaultLogLevel     = "info"
    DefaultMCPTransport = "stdio"
    DefaultMCPHTTPPort  = "8080"
)

// Config holds all application configuration
type Config struct {
    PaperlessURL   string
    PaperlessToken string
    MCPAuthToken   string // optional
    LogLevel       string
    MCPTransport   string
    MCPHTTPPort    string
}

// Load reads configuration from environment variables
func Load() (*Config, error) {
    cfg := &Config{}

    cfg.PaperlessURL = os.Getenv(EnvPaperlessURL)
    if strings.TrimSpace(cfg.PaperlessURL) == "" {
        return nil, errors.New("environment variable PAPERLESS_URL is required but not set")
    }

    cfg.PaperlessToken = os.Getenv(EnvPaperlessToken)
    if strings.TrimSpace(cfg.PaperlessToken) == "" {
        return nil, errors.New("environment variable PAPERLESS_TOKEN is required but not set")
    }

    cfg.MCPAuthToken = os.Getenv(EnvMCPAuthToken) // optional, no error if empty

    // Optional vars with defaults
    cfg.LogLevel = os.Getenv(EnvLogLevel)
    if cfg.LogLevel == "" {
        cfg.LogLevel = DefaultLogLevel
    }
    cfg.LogLevel = strings.ToLower(cfg.LogLevel)
    allowedLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
    if !allowedLogLevels[cfg.LogLevel] {
        return nil, fmt.Errorf("invalid log level: %s, allowed: debug, info, warn, error", cfg.LogLevel)
    }

    cfg.MCPTransport = os.Getenv(EnvMCPTransport)
    if cfg.MCPTransport == "" {
        cfg.MCPTransport = DefaultMCPTransport
    }
    cfg.MCPTransport = strings.ToLower(cfg.MCPTransport)
    if cfg.MCPTransport != "stdio" && cfg.MCPTransport != "http" {
        return nil, fmt.Errorf("invalid MCP_TRANSPORT: %s, allowed: stdio, http", cfg.MCPTransport)
    }

    cfg.MCPHTTPPort = os.Getenv(EnvMCPHTTPPort)
    if cfg.MCPHTTPPort == "" {
        cfg.MCPHTTPPort = DefaultMCPHTTPPort
    }
    // Optional: Could add port format validation here but skipping per spec simplicity

    return cfg, nil
}
