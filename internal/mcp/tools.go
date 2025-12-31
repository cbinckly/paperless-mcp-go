package mcp

import (
	"context"
	"log/slog"
)

// registerTools registers all MCP tools with the server
func (s *Server) registerTools() {
	slog.Debug("Registering MCP tools")

	// Register the ping test tool
	err := s.RegisterTool(Tool{
		Name:        "ping",
		Description: "Test tool that returns pong - useful for verifying the MCP server is working",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
			"required":   []string{},
		},
		Handler: s.handlePing,
	})
	if err != nil {
		slog.Error("Failed to register ping tool", "error", err)
	}

	// Register the server_info tool
	err = s.RegisterTool(Tool{
		Name:        "server_info",
		Description: "Returns information about the MCP server and Paperless connection",
		InputSchema: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
			"required":   []string{},
		},
		Handler: s.handleServerInfo,
	})
	if err != nil {
		slog.Error("Failed to register server_info tool", "error", err)
	}


	// Register the search_documents tool
	err = s.RegisterTool(Tool{
		Name:        "search_documents",
		Description: "Search for documents in Paperless by text query with pagination support",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{
					"type":        "string",
					"description": "Search query text",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number (1-based, optional, default: 1)",
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "Number of results per page (optional, default: 25, max: 100)",
				},
			},
			"required": []string{"query"},
		},
		Handler: s.handleSearchDocuments,
	})
	if err != nil {
		slog.Error("Failed to register search_documents tool", "error", err)
	}

	// Register the find_similar_documents tool
	err = s.RegisterTool(Tool{
		Name:        "find_similar_documents",
		Description: "Find documents similar to a given document with pagination support",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the document to find similar documents for",
				},
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number (1-based, optional, default: 1)",
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "Number of results per page (optional, default: 25, max: 100)",
				},
			},
			"required": []string{"document_id"},
		},
		Handler: s.handleFindSimilarDocuments,
	})
	if err != nil {
		slog.Error("Failed to register find_similar_documents tool", "error", err)
	}

	slog.Info("Tool registration complete", "total_tools", len(s.tools))
}

// handlePing is a simple test tool that returns "pong"
func (s *Server) handlePing(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	slog.Debug("Ping tool invoked")
	return map[string]string{
		"status":  "ok",
		"message": "pong",
	}, nil
}

// handleServerInfo returns information about the MCP server
func (s *Server) handleServerInfo(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	slog.Debug("Server info tool invoked")
	return map[string]string{
		"server_name":    ServerName,
		"server_version": ServerVersion,
		"paperless_url":  s.cfg.PaperlessURL,
		"transport":      s.cfg.MCPTransport,
		"status":         "ok",
	}, nil
}
