package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// Default pagination values for document operations
const (
	DefaultPage     = 1
	DefaultPageSize = 25
	MaxPageSize     = 100
)

// handleSearchDocuments handles the search_documents tool
func (s *Server) handleSearchDocuments(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate query
	query, ok := args["query"].(string)
	if !ok || query == "" {
		return nil, fmt.Errorf("query parameter is required and must be a non-empty string")
	}

	// Extract optional page parameter
	page := DefaultPage
	if pageVal, ok := args["page"].(float64); ok {
		page = int(pageVal)
		if page < 1 {
			page = DefaultPage
		}
	}

	// Extract optional page_size parameter
	pageSize := DefaultPageSize
	if pageSizeVal, ok := args["page_size"].(float64); ok {
		pageSize = int(pageSizeVal)
		if pageSize < 1 {
			pageSize = DefaultPageSize
		} else if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
	}

	slog.Debug("Searching documents",
		"query", query,
		"page", page,
		"page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.SearchDocuments(ctx, query, page, pageSize)
	if err != nil {
		slog.Error("Failed to search documents",
			"query", query,
			"error", err)
		return nil, fmt.Errorf("failed to search documents: %w", err)
	}

	// Parse documents from Results
	var documents []paperless.Document
	if err := json.Unmarshal(response.Results, &documents); err != nil {
		slog.Error("Failed to parse search results", "error", err)
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	slog.Info("Documents search completed",
		"query", query,
		"found", response.Count,
		"returned", len(documents))

	return map[string]interface{}{
		"count":      response.Count,
		"page":       page,
		"page_size":  pageSize,
		"has_next":   response.Next != nil,
		"has_prev":   response.Previous != nil,
		"documents":  documents,
	}, nil
}

// handleFindSimilarDocuments handles the find_similar_documents tool
func (s *Server) handleFindSimilarDocuments(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_id
	documentIDFloat, ok := args["document_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_id parameter is required and must be an integer")
	}
	documentID := int(documentIDFloat)
	if documentID < 1 {
		return nil, fmt.Errorf("document_id must be a positive integer")
	}

	// Extract optional page parameter
	page := DefaultPage
	if pageVal, ok := args["page"].(float64); ok {
		page = int(pageVal)
		if page < 1 {
			page = DefaultPage
		}
	}

	// Extract optional page_size parameter
	pageSize := DefaultPageSize
	if pageSizeVal, ok := args["page_size"].(float64); ok {
		pageSize = int(pageSizeVal)
		if pageSize < 1 {
			pageSize = DefaultPageSize
		} else if pageSize > MaxPageSize {
			pageSize = MaxPageSize
		}
	}

	slog.Debug("Finding similar documents",
		"document_id", documentID,
		"page", page,
		"page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.GetSimilarDocuments(ctx, documentID, page, pageSize)
	if err != nil {
		slog.Error("Failed to find similar documents",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to find similar documents: %w", err)
	}

	// Parse documents from Results
	var documents []paperless.Document
	if err := json.Unmarshal(response.Results, &documents); err != nil {
		slog.Error("Failed to parse similar documents results", "error", err)
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	slog.Info("Similar documents search completed",
		"document_id", documentID,
		"found", response.Count,
		"returned", len(documents))

	return map[string]interface{}{
		"document_id": documentID,
		"count":       response.Count,
		"page":        page,
		"page_size":   pageSize,
		"has_next":    response.Next != nil,
		"has_prev":    response.Previous != nil,
		"documents":   documents,
	}, nil
}
