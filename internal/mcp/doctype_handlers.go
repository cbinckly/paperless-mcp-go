package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// handleListDocumentTypes handles the list_document_types tool
func (s *Server) handleListDocumentTypes(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	slog.Debug("Listing document types", "page", page, "page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.ListDocumentTypes(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to list document types", "error", err)
		return nil, fmt.Errorf("failed to list document types: %w", err)
	}

	// Parse document types from Results
	var documentTypes []paperless.DocumentType
	if err := json.Unmarshal(response.Results, &documentTypes); err != nil {
		slog.Error("Failed to parse document types results", "error", err)
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	slog.Info("Document types listed successfully",
		"count", response.Count,
		"returned", len(documentTypes))

	return map[string]interface{}{
		"count":          response.Count,
		"page":           page,
		"page_size":      pageSize,
		"has_next":       response.Next != nil,
		"has_prev":       response.Previous != nil,
		"document_types": documentTypes,
	}, nil
}

// handleGetDocumentType handles the get_document_type tool
func (s *Server) handleGetDocumentType(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_type_id
	documentTypeIDFloat, ok := args["document_type_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_type_id parameter is required and must be an integer")
	}
	documentTypeID := int(documentTypeIDFloat)
	if documentTypeID < 1 {
		return nil, fmt.Errorf("document_type_id must be a positive integer")
	}

	slog.Debug("Getting document type", "document_type_id", documentTypeID)

	// Call Paperless API
	documentType, err := s.paperlessClient.GetDocumentType(ctx, documentTypeID)
	if err != nil {
		slog.Error("Failed to get document type",
			"document_type_id", documentTypeID,
			"error", err)
		return nil, fmt.Errorf("failed to get document type: %w", err)
	}

	slog.Info("Document type retrieved successfully",
		"document_type_id", documentTypeID,
		"name", documentType.Name)

	return documentType, nil
}

// handleCreateDocumentType handles the create_document_type tool
func (s *Server) handleCreateDocumentType(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract required name
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required and must be a non-empty string")
	}

	slog.Debug("Creating document type", "name", name)

	// Build document type from args
	documentType := &paperless.DocumentType{
		Name: name,
	}

	// Extract optional fields
	if match, ok := args["match"].(string); ok {
		documentType.Match = match
	}
	if matchingAlg, ok := args["matching_algorithm"].(float64); ok {
		documentType.MatchingAlgorithm = int(matchingAlg)
	}
	if isInsensitive, ok := args["is_insensitive"].(bool); ok {
		documentType.IsInsensitive = isInsensitive
	}

	// Call Paperless API
	createdDocumentType, err := s.paperlessClient.CreateDocumentType(ctx, documentType)
	if err != nil {
		slog.Error("Failed to create document type",
			"name", name,
			"error", err)
		return nil, fmt.Errorf("failed to create document type: %w", err)
	}

	slog.Info("Document type created successfully",
		"document_type_id", createdDocumentType.ID,
		"name", createdDocumentType.Name)

	return createdDocumentType, nil
}

// handleUpdateDocumentType handles the update_document_type tool
func (s *Server) handleUpdateDocumentType(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_type_id
	documentTypeIDFloat, ok := args["document_type_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_type_id parameter is required and must be an integer")
	}
	documentTypeID := int(documentTypeIDFloat)
	if documentTypeID < 1 {
		return nil, fmt.Errorf("document_type_id must be a positive integer")
	}

	// Build updates map from args (exclude document_type_id)
	updates := make(map[string]interface{})
	for key, value := range args {
		if key != "document_type_id" {
			updates[key] = value
		}
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field to update must be provided")
	}

	slog.Debug("Updating document type",
		"document_type_id", documentTypeID,
		"fields", len(updates))

	// Call Paperless API
	updatedDocumentType, err := s.paperlessClient.UpdateDocumentType(ctx, documentTypeID, updates)
	if err != nil {
		slog.Error("Failed to update document type",
			"document_type_id", documentTypeID,
			"error", err)
		return nil, fmt.Errorf("failed to update document type: %w", err)
	}

	slog.Info("Document type updated successfully",
		"document_type_id", documentTypeID,
		"name", updatedDocumentType.Name)

	return updatedDocumentType, nil
}

// handleDeleteDocumentType handles the delete_document_type tool
func (s *Server) handleDeleteDocumentType(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_type_id
	documentTypeIDFloat, ok := args["document_type_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_type_id parameter is required and must be an integer")
	}
	documentTypeID := int(documentTypeIDFloat)
	if documentTypeID < 1 {
		return nil, fmt.Errorf("document_type_id must be a positive integer")
	}

	slog.Debug("Deleting document type", "document_type_id", documentTypeID)

	// Call Paperless API
	err := s.paperlessClient.DeleteDocumentType(ctx, documentTypeID)
	if err != nil {
		slog.Error("Failed to delete document type",
			"document_type_id", documentTypeID,
			"error", err)
		return nil, fmt.Errorf("failed to delete document type: %w", err)
	}

	slog.Info("Document type deleted successfully", "document_type_id", documentTypeID)

	return map[string]interface{}{
		"success":          true,
		"document_type_id": documentTypeID,
		"message":          "Document type deleted successfully",
	}, nil
}
