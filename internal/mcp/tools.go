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


	// Register the get_document tool
	err = s.RegisterTool(Tool{
		Name:        "get_document",
		Description: "Get a document by ID with all metadata",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the document to retrieve",
				},
			},
			"required": []string{"document_id"},
		},
		Handler: s.handleGetDocument,
	})
	if err != nil {
		slog.Error("Failed to register get_document tool", "error", err)
	}

	// Register the get_document_content tool
	err = s.RegisterTool(Tool{
		Name:        "get_document_content",
		Description: "Get the text content of a document",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the document to retrieve content from",
				},
			},
			"required": []string{"document_id"},
		},
		Handler: s.handleGetDocumentContent,
	})
	if err != nil {
		slog.Error("Failed to register get_document_content tool", "error", err)
	}

	// Register the create_document tool
	err = s.RegisterTool(Tool{
		Name:        "create_document",
		Description: "Create a new document in Paperless",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"title": map[string]interface{}{
					"type":        "string",
					"description": "Title of the document",
				},
				"correspondent": map[string]interface{}{
					"type":        "integer",
					"description": "Correspondent ID (optional)",
				},
				"document_type": map[string]interface{}{
					"type":        "integer",
					"description": "Document type ID (optional)",
				},
				"storage_path": map[string]interface{}{
					"type":        "integer",
					"description": "Storage path ID (optional)",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"description": "Array of tag IDs (optional)",
					"items": map[string]interface{}{
						"type": "integer",
					},
				},
			},
			"required": []string{"title"},
		},
		Handler: s.handleCreateDocument,
	})
	if err != nil {
		slog.Error("Failed to register create_document tool", "error", err)
	}

	// Register the update_document tool
	err = s.RegisterTool(Tool{
		Name:        "update_document",
		Description: "Update a document's metadata",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the document to update",
				},
				"title": map[string]interface{}{
					"type":        "string",
					"description": "New title (optional)",
				},
				"correspondent": map[string]interface{}{
					"type":        "integer",
					"description": "New correspondent ID (optional)",
				},
				"document_type": map[string]interface{}{
					"type":        "integer",
					"description": "New document type ID (optional)",
				},
				"storage_path": map[string]interface{}{
					"type":        "integer",
					"description": "New storage path ID (optional)",
				},
				"tags": map[string]interface{}{
					"type":        "array",
					"description": "New array of tag IDs (optional)",
					"items": map[string]interface{}{
						"type": "integer",
					},
				},
			},
			"required": []string{"document_id"},
		},
		Handler: s.handleUpdateDocument,
	})
	if err != nil {
		slog.Error("Failed to register update_document tool", "error", err)
	}

	// Register the delete_document tool
	err = s.RegisterTool(Tool{
		Name:        "delete_document",
		Description: "Delete a document from Paperless",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"document_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the document to delete",
				},
			},
			"required": []string{"document_id"},
		},
		Handler: s.handleDeleteDocument,
	})
	if err != nil {
		slog.Error("Failed to register delete_document tool", "error", err)
	}


	// Register the list_correspondents tool
	err = s.RegisterTool(Tool{
		Name:        "list_correspondents",
		Description: "List all correspondents with pagination support",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"page": map[string]interface{}{
					"type":        "integer",
					"description": "Page number (1-based, optional, default: 1)",
				},
				"page_size": map[string]interface{}{
					"type":        "integer",
					"description": "Number of results per page (optional, default: 25, max: 100)",
				},
			},
			"required": []string{},
		},
		Handler: s.handleListCorrespondents,
	})
	if err != nil {
		slog.Error("Failed to register list_correspondents tool", "error", err)
	}

	// Register the get_correspondent tool
	err = s.RegisterTool(Tool{
		Name:        "get_correspondent",
		Description: "Get a correspondent by ID",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"correspondent_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the correspondent to retrieve",
				},
			},
			"required": []string{"correspondent_id"},
		},
		Handler: s.handleGetCorrespondent,
	})
	if err != nil {
		slog.Error("Failed to register get_correspondent tool", "error", err)
	}

	// Register the create_correspondent tool
	err = s.RegisterTool(Tool{
		Name:        "create_correspondent",
		Description: "Create a new correspondent in Paperless",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"name": map[string]interface{}{
					"type":        "string",
					"description": "Name of the correspondent",
				},
				"match": map[string]interface{}{
					"type":        "string",
					"description": "Matching text pattern (optional)",
				},
				"matching_algorithm": map[string]interface{}{
					"type":        "integer",
					"description": "Matching algorithm type (optional)",
				},
				"is_insensitive": map[string]interface{}{
					"type":        "boolean",
					"description": "Case insensitive matching (optional)",
				},
			},
			"required": []string{"name"},
		},
		Handler: s.handleCreateCorrespondent,
	})
	if err != nil {
		slog.Error("Failed to register create_correspondent tool", "error", err)
	}

	// Register the update_correspondent tool
	err = s.RegisterTool(Tool{
		Name:        "update_correspondent",
		Description: "Update a correspondent's information",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"correspondent_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the correspondent to update",
				},
				"name": map[string]interface{}{
					"type":        "string",
					"description": "New name (optional)",
				},
				"match": map[string]interface{}{
					"type":        "string",
					"description": "New matching pattern (optional)",
				},
				"matching_algorithm": map[string]interface{}{
					"type":        "integer",
					"description": "New matching algorithm (optional)",
				},
				"is_insensitive": map[string]interface{}{
					"type":        "boolean",
					"description": "Case insensitive matching (optional)",
				},
			},
			"required": []string{"correspondent_id"},
		},
		Handler: s.handleUpdateCorrespondent,
	})
	if err != nil {
		slog.Error("Failed to register update_correspondent tool", "error", err)
	}

	// Register the delete_correspondent tool
	err = s.RegisterTool(Tool{
		Name:        "delete_correspondent",
		Description: "Delete a correspondent from Paperless",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"correspondent_id": map[string]interface{}{
					"type":        "integer",
					"description": "ID of the correspondent to delete",
				},
			},
			"required": []string{"correspondent_id"},
		},
		Handler: s.handleDeleteCorrespondent,
	})
	if err != nil {
		slog.Error("Failed to register delete_correspondent tool", "error", err)
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
