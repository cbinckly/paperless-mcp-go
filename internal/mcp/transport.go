package mcp

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mark3labs/mcp-go/server"
)

// Transport constants
const (
	// ShutdownTimeout is the maximum time to wait for graceful shutdown
	ShutdownTimeout = 10 * time.Second

	// StreamableHTTPEndpoint is the HTTP streaming endpoint path
	StreamableHTTPEndpoint = "/mcp"

	// HeartbeatInterval is the interval for connection heartbeat/keep-alive
	HeartbeatInterval = 30 * time.Second

	// HealthEndpoint is the health check endpoint
	HealthEndpoint = "/health"

	// HTTP server timeout constants
	HTTPReadTimeout  = 30 * time.Second
	HTTPWriteTimeout = 0                  // No timeout for streaming
	HTTPIdleTimeout  = 120 * time.Second
)

// StartStdio starts the MCP server with stdio transport
func (s *Server) StartStdio(ctx context.Context) error {
	slog.Info("Starting MCP server with stdio transport")

	// Create a channel to listen for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create error channel
	errChan := make(chan error, 1)

	// Start the stdio server in a goroutine
	go func() {
		slog.Debug("Starting stdio transport listener")
		// Use the mark3labs MCP SDK stdio server
		if err := server.ServeStdio(s.mcpServer); err != nil {
			slog.Error("Stdio server error", "error", err)
			errChan <- err
		}
	}()

	// Wait for shutdown signal or error
	select {
	case <-ctx.Done():
		slog.Info("Context cancelled, shutting down stdio server")
		return nil
	case sig := <-sigChan:
		slog.Info("Received shutdown signal", "signal", sig)
		return nil
	case err := <-errChan:
		return fmt.Errorf("stdio server error: %w", err)
	}
}

// StartHTTP starts the MCP server with StreamableHTTP transport
func (s *Server) StartHTTP(ctx context.Context) error {
	port := s.cfg.MCPHTTPPort
	addr := ":" + port
	slog.Info("Starting MCP server with StreamableHTTP transport",
		"port", port,
		"endpoint", StreamableHTTPEndpoint,
		"heartbeat_interval", HeartbeatInterval)

	// Create the StreamableHTTP server from the SDK
	// This replaces the legacy SSE transport with the modern HTTP streaming standard
	streamableServer := server.NewStreamableHTTPServer(s.mcpServer,
		server.WithEndpointPath(StreamableHTTPEndpoint),
		server.WithHeartbeatInterval(HeartbeatInterval),
	)

	// Create HTTP server mux
	mux := http.NewServeMux()

	// Setup health endpoint
	mux.HandleFunc(HealthEndpoint, s.handleHealth)

	// Setup StreamableHTTP endpoint using the SDK's server
	// StreamableHTTP handles POST (client messages), GET (server notifications), and DELETE (cleanup)
	mux.Handle(StreamableHTTPEndpoint, streamableServer)

	// Create HTTP server with timeouts
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      s.authMiddleware(mux),
		ReadTimeout:  HTTPReadTimeout,
		WriteTimeout: HTTPWriteTimeout,
		IdleTimeout:  HTTPIdleTimeout,
	}

	// Create a channel to listen for shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Create error channel
	errChan := make(chan error, 1)

	// Start the HTTP server in a goroutine
	go func() {
		slog.Info("HTTP server listening", "addr", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server error", "error", err)
			errChan <- err
		}
	}()

	// Wait for shutdown signal, context cancellation, or error
	select {
	case <-ctx.Done():
		slog.Info("Context cancelled, initiating HTTP server shutdown")
	case sig := <-sigChan:
		slog.Info("Received shutdown signal, initiating HTTP server shutdown", "signal", sig)
	case err := <-errChan:
		return fmt.Errorf("HTTP server error: %w", err)
	}

	// Graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	slog.Info("Shutting down HTTP server gracefully", "timeout", ShutdownTimeout)
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown error", "error", err)
		return fmt.Errorf("shutdown error: %w", err)
	}

	slog.Info("HTTP server shutdown complete")
	return nil
}

// authMiddleware adds authentication if MCP_AUTH_TOKEN is configured
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If no auth token is configured, skip authentication
		if s.cfg.MCPAuthToken == "" {
			next.ServeHTTP(w, r)
			return
		}

		// Allow health check without authentication
		if r.URL.Path == HealthEndpoint {
			next.ServeHTTP(w, r)
			return
		}

		// Check Authorization header
		authHeader := r.Header.Get("Authorization")
		expectedAuth := "Bearer " + s.cfg.MCPAuthToken

		if authHeader != expectedAuth {
			slog.Warn("Authentication failed",
				"path", r.URL.Path,
				"remote_addr", r.RemoteAddr)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		slog.Debug("Authentication successful",
			"path", r.URL.Path,
			"remote_addr", r.RemoteAddr)

		next.ServeHTTP(w, r)
	})
}

// handleHealth handles the health check endpoint
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok","server":"` + ServerName + `","version":"` + ServerVersion + `"}`))
}
