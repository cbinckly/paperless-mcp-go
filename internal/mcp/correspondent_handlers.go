package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// handleListCorrespondents handles the list_correspondents tool
func (s *Server) handleListCorrespondents(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	slog.Debug("Listing correspondents", "page", page, "page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.ListCorrespondents(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to list correspondents", "error", err)
		return nil, fmt.Errorf("failed to list correspondents: %w", err)
	}

	// Parse correspondents from Results
	var correspondents []paperless.Correspondent
	if err := json.Unmarshal(response.Results, &correspondents); err != nil {
		slog.Error("Failed to parse correspondents results", "error", err)
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	slog.Info("Correspondents listed successfully",
		"count", response.Count,
		"returned", len(correspondents))

	return map[string]interface{}{
		"count":          response.Count,
		"page":           page,
		"page_size":      pageSize,
		"has_next":       response.Next != nil,
		"has_prev":       response.Previous != nil,
		"correspondents": correspondents,
	}, nil
}

// handleGetCorrespondent handles the get_correspondent tool
func (s *Server) handleGetCorrespondent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate correspondent_id
	correspondentIDFloat, ok := args["correspondent_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("correspondent_id parameter is required and must be an integer")
	}
	correspondentID := int(correspondentIDFloat)
	if correspondentID < 1 {
		return nil, fmt.Errorf("correspondent_id must be a positive integer")
	}

	slog.Debug("Getting correspondent", "correspondent_id", correspondentID)

	// Call Paperless API
	correspondent, err := s.paperlessClient.GetCorrespondent(ctx, correspondentID)
	if err != nil {
		slog.Error("Failed to get correspondent",
			"correspondent_id", correspondentID,
			"error", err)
		return nil, fmt.Errorf("failed to get correspondent: %w", err)
	}

	slog.Info("Correspondent retrieved successfully",
		"correspondent_id", correspondentID,
		"name", correspondent.Name)

	return correspondent, nil
}

// handleCreateCorrespondent handles the create_correspondent tool
func (s *Server) handleCreateCorrespondent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract required name
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required and must be a non-empty string")
	}

	slog.Debug("Creating correspondent", "name", name)

	// Build correspondent from args
	correspondent := &paperless.Correspondent{
		Name: name,
	}

	// Extract optional fields
	if match, ok := args["match"].(string); ok {
		correspondent.Match = match
	}
	if matchingAlg, ok := args["matching_algorithm"].(float64); ok {
		correspondent.MatchingAlgorithm = int(matchingAlg)
	}
	if isInsensitive, ok := args["is_insensitive"].(bool); ok {
		correspondent.IsInsensitive = isInsensitive
	}

	// Call Paperless API
	createdCorrespondent, err := s.paperlessClient.CreateCorrespondent(ctx, correspondent)
	if err != nil {
		slog.Error("Failed to create correspondent",
			"name", name,
			"error", err)
		return nil, fmt.Errorf("failed to create correspondent: %w", err)
	}

	slog.Info("Correspondent created successfully",
		"correspondent_id", createdCorrespondent.ID,
		"name", createdCorrespondent.Name)

	return createdCorrespondent, nil
}

// handleUpdateCorrespondent handles the update_correspondent tool
func (s *Server) handleUpdateCorrespondent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate correspondent_id
	correspondentIDFloat, ok := args["correspondent_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("correspondent_id parameter is required and must be an integer")
	}
	correspondentID := int(correspondentIDFloat)
	if correspondentID < 1 {
		return nil, fmt.Errorf("correspondent_id must be a positive integer")
	}

	// Build updates map from args (exclude correspondent_id)
	updates := make(map[string]interface{})
	for key, value := range args {
		if key != "correspondent_id" {
			updates[key] = value
		}
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field to update must be provided")
	}

	slog.Debug("Updating correspondent",
		"correspondent_id", correspondentID,
		"fields", len(updates))

	// Call Paperless API
	updatedCorrespondent, err := s.paperlessClient.UpdateCorrespondent(ctx, correspondentID, updates)
	if err != nil {
		slog.Error("Failed to update correspondent",
			"correspondent_id", correspondentID,
			"error", err)
		return nil, fmt.Errorf("failed to update correspondent: %w", err)
	}

	slog.Info("Correspondent updated successfully",
		"correspondent_id", correspondentID,
		"name", updatedCorrespondent.Name)

	return updatedCorrespondent, nil
}

// handleDeleteCorrespondent handles the delete_correspondent tool
func (s *Server) handleDeleteCorrespondent(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate correspondent_id
	correspondentIDFloat, ok := args["correspondent_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("correspondent_id parameter is required and must be an integer")
	}
	correspondentID := int(correspondentIDFloat)
	if correspondentID < 1 {
		return nil, fmt.Errorf("correspondent_id must be a positive integer")
	}

	slog.Debug("Deleting correspondent", "correspondent_id", correspondentID)

	// Call Paperless API
	err := s.paperlessClient.DeleteCorrespondent(ctx, correspondentID)
	if err != nil {
		slog.Error("Failed to delete correspondent",
			"correspondent_id", correspondentID,
			"error", err)
		return nil, fmt.Errorf("failed to delete correspondent: %w", err)
	}

	slog.Info("Correspondent deleted successfully", "correspondent_id", correspondentID)

	return map[string]interface{}{
		"success":          true,
		"correspondent_id": correspondentID,
		"message":          "Correspondent deleted successfully",
	}, nil
}
