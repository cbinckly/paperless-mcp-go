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

// handleGetDocument handles the get_document tool
func (s *Server) handleGetDocument(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_id
	documentIDFloat, ok := args["document_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_id parameter is required and must be an integer")
	}
	documentID := int(documentIDFloat)
	if documentID < 1 {
		return nil, fmt.Errorf("document_id must be a positive integer")
	}

	slog.Debug("Getting document", "document_id", documentID)

	// Call Paperless API
	document, err := s.paperlessClient.GetDocument(ctx, documentID)
	if err != nil {
		slog.Error("Failed to get document",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	slog.Info("Document retrieved successfully",
		"document_id", documentID,
		"title", document.Title)

	return document, nil
}

// handleGetDocumentContent handles the get_document_content tool
func (s *Server) handleGetDocumentContent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_id
	documentIDFloat, ok := args["document_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_id parameter is required and must be an integer")
	}
	documentID := int(documentIDFloat)
	if documentID < 1 {
		return nil, fmt.Errorf("document_id must be a positive integer")
	}

	slog.Debug("Getting document content", "document_id", documentID)

	// Call Paperless API
	content, err := s.paperlessClient.GetDocumentContent(ctx, documentID)
	if err != nil {
		slog.Error("Failed to get document content",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to get document content: %w", err)
	}

	slog.Info("Document content retrieved successfully",
		"document_id", documentID,
		"content_length", len(content))

	return map[string]interface{}{
		"document_id": documentID,
		"content":     content,
	}, nil
}

// handleCreateDocument handles the create_document tool
func (s *Server) handleCreateDocument(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract required title
	title, ok := args["title"].(string)
	if !ok || title == "" {
		return nil, fmt.Errorf("title parameter is required and must be a non-empty string")
	}

	slog.Debug("Creating document", "title", title)

	// Build document from args
	document := &paperless.Document{
		Title: title,
	}

	// Extract optional fields
	if correspondent, ok := args["correspondent"].(float64); ok {
		correspondentID := int(correspondent)
		document.Correspondent = &correspondentID
	}
	if docType, ok := args["document_type"].(float64); ok {
		docTypeID := int(docType)
		document.DocumentType = &docTypeID
	}
	if storagePath, ok := args["storage_path"].(float64); ok {
		storagePathID := int(storagePath)
		document.StoragePath = &storagePathID
	}
	if tagsInterface, ok := args["tags"].([]interface{}); ok {
		tags := make([]int, 0, len(tagsInterface))
		for _, tagInterface := range tagsInterface {
			if tagFloat, ok := tagInterface.(float64); ok {
				tags = append(tags, int(tagFloat))
			}
		}
		document.Tags = tags
	}

	// Call Paperless API
	createdDocument, err := s.paperlessClient.CreateDocument(ctx, document)
	if err != nil {
		slog.Error("Failed to create document",
			"title", title,
			"error", err)
		return nil, fmt.Errorf("failed to create document: %w", err)
	}

	slog.Info("Document created successfully",
		"document_id", createdDocument.ID,
		"title", createdDocument.Title)

	return createdDocument, nil
}

// handleUpdateDocument handles the update_document tool
func (s *Server) handleUpdateDocument(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_id
	documentIDFloat, ok := args["document_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_id parameter is required and must be an integer")
	}
	documentID := int(documentIDFloat)
	if documentID < 1 {
		return nil, fmt.Errorf("document_id must be a positive integer")
	}

	// Build updates map from args (exclude document_id)
	updates := make(map[string]interface{})
	for key, value := range args {
		if key != "document_id" {
			updates[key] = value
		}
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field to update must be provided")
	}

	slog.Debug("Updating document",
		"document_id", documentID,
		"fields", len(updates))

	// Call Paperless API
	updatedDocument, err := s.paperlessClient.UpdateDocument(ctx, documentID, updates)
	if err != nil {
		slog.Error("Failed to update document",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to update document: %w", err)
	}

	slog.Info("Document updated successfully",
		"document_id", documentID,
		"title", updatedDocument.Title)

	return updatedDocument, nil
}

// handleDeleteDocument handles the delete_document tool
func (s *Server) handleDeleteDocument(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate document_id
	documentIDFloat, ok := args["document_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("document_id parameter is required and must be an integer")
	}
	documentID := int(documentIDFloat)
	if documentID < 1 {
		return nil, fmt.Errorf("document_id must be a positive integer")
	}

	slog.Debug("Deleting document", "document_id", documentID)

	// Call Paperless API
	err := s.paperlessClient.DeleteDocument(ctx, documentID)
	if err != nil {
		slog.Error("Failed to delete document",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to delete document: %w", err)
	}

	slog.Info("Document deleted successfully", "document_id", documentID)

	return map[string]interface{}{
		"success":     true,
		"document_id": documentID,
		"message":     "Document deleted successfully",
	}, nil
}
