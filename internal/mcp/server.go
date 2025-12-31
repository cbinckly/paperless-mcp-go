// Package mcp provides the Model Context Protocol server implementation
// for the Paperless MCP server with dual transport support (stdio and HTTP/SSE).
package mcp

import (
	"context"
	"encoding/json"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/config"
	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Server name and version constants
const (
	ServerName    = "Paperless MCP Server"
	ServerVersion = "1.0.0"
)

// Server represents the MCP server
type Server struct {
	cfg             *config.Config
	paperlessClient *paperless.Client
	mcpServer       *server.MCPServer
	tools           map[string]Tool
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string
	Description string
	InputSchema map[string]interface{}
	Handler     ToolHandler
}

// ToolHandler is the function signature for tool handlers
type ToolHandler func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// New creates a new MCP server instance
func New(cfg *config.Config) (*Server, error) {
	slog.Debug("Creating new MCP server",
		"paperless_url", cfg.PaperlessURL,
		"transport", cfg.MCPTransport)

	// Create Paperless client
	paperlessClient := paperless.New(cfg.PaperlessURL, cfg.PaperlessToken)

	// Create MCP server instance with the mark3labs SDK
	mcpServer := server.NewMCPServer(
		ServerName,
		ServerVersion,
		server.WithLogging(),
	)

	s := &Server{
		cfg:             cfg,
		paperlessClient: paperlessClient,
		mcpServer:       mcpServer,
		tools:           make(map[string]Tool),
	}

	// Register initial tools
	s.registerTools()

	slog.Info("MCP server created successfully",
		"server_name", ServerName,
		"server_version", ServerVersion,
		"tool_count", len(s.tools))

	return s, nil
}

// RegisterTool registers a new tool with the MCP server
func (s *Server) RegisterTool(tool Tool) error {
	slog.Debug("Registering tool",
		"tool_name", tool.Name,
		"description", tool.Description)

	// Store in our tools map
	s.tools[tool.Name] = tool

	// Create the MCP tool using the SDK with just name and description
	// The schema will be handled by the SDK
	mcpTool := mcp.NewTool(tool.Name,
		mcp.WithDescription(tool.Description),
	)

	// Create the handler wrapper that calls our tool handler
	toolName := tool.Name // Capture for closure
	handlerWrapper := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		// Extract arguments from the request with proper type assertion
		args := make(map[string]interface{})
		if request.Params.Arguments != nil {
			// Try type assertion first
			if argsMap, ok := request.Params.Arguments.(map[string]interface{}); ok {
				args = argsMap
			} else {
				// If direct assertion fails, try JSON round-trip
				jsonBytes, err := json.Marshal(request.Params.Arguments)
				if err == nil {
					json.Unmarshal(jsonBytes, &args)
				}
			}
		}

		// Call our tool handler
		result, err := s.ExecuteTool(ctx, toolName, args)
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		// Convert result to string representation
		resultStr := formatToolResult(result)
		return mcp.NewToolResultText(resultStr), nil
	}

	// Add the tool to the MCP server
	s.mcpServer.AddTool(mcpTool, handlerWrapper)

	slog.Info("Tool registered successfully", "tool_name", tool.Name)
	return nil
}

// formatToolResult formats a tool result for MCP response
func formatToolResult(result interface{}) string {
	// Try JSON encoding first for structured data
	jsonBytes, err := json.Marshal(result)
	if err == nil {
		return string(jsonBytes)
	}
	
	// Fallback to simple string formatting
	switch v := result.(type) {
	case string:
		return v
	case map[string]interface{}:
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	case map[string]string:
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	default:
		jsonBytes, _ := json.Marshal(v)
		return string(jsonBytes)
	}
}

// GetPaperlessClient returns the Paperless API client
func (s *Server) GetPaperlessClient() *paperless.Client {
	return s.paperlessClient
}

// GetConfig returns the server configuration
func (s *Server) GetConfig() *config.Config {
	return s.cfg
}
