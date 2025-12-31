package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// handleListTags handles the list_tags tool
func (s *Server) handleListTags(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract pagination parameters
	page := 1
	pageSize := paperless.DefaultPageSize

	if p, ok := args["page"].(float64); ok {
		page = int(p)
	}
	if ps, ok := args["page_size"].(float64); ok {
		pageSize = int(ps)
	}

	slog.Debug("List tags tool invoked", "page", page, "page_size", pageSize)

	// Call API
	response, err := s.paperlessClient.ListTags(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to list tags", "error", err)
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	// Parse results as tags
	var tags []paperless.Tag
	if err := json.Unmarshal(response.Results, &tags); err != nil {
		slog.Error("Failed to parse tags from response", "error", err)
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	return map[string]interface{}{
		"count":    response.Count,
		"next":     response.Next,
		"previous": response.Previous,
		"tags":     tags,
	}, nil
}

// handleGetTag handles the get_tag tool
func (s *Server) handleGetTag(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract tag ID
	tagID, ok := args["tag_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("tag_id is required and must be an integer")
	}

	slog.Debug("Get tag tool invoked", "tag_id", int(tagID))

	// Call API
	tag, err := s.paperlessClient.GetTag(ctx, int(tagID))
	if err != nil {
		slog.Error("Failed to get tag", "tag_id", int(tagID), "error", err)
		return nil, fmt.Errorf("failed to get tag: %w", err)
	}

	return tag, nil
}

// handleCreateTag handles the create_tag tool
func (s *Server) handleCreateTag(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate required fields
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name is required and must be a non-empty string")
	}

	color, ok := args["color"].(string)
	if !ok || color == "" {
		return nil, fmt.Errorf("color is required and must be a non-empty string")
	}

	slog.Debug("Create tag tool invoked", "name", name, "color", color)

	// Build tag object
	tag := &paperless.Tag{
		Name:  name,
		Color: color,
	}

	// Optional fields
	if match, ok := args["match"].(string); ok {
		tag.Match = match
	}
	if matchingAlgo, ok := args["matching_algorithm"].(float64); ok {
		tag.MatchingAlgorithm = int(matchingAlgo)
	}
	if isInsensitive, ok := args["is_insensitive"].(bool); ok {
		tag.IsInsensitive = isInsensitive
	}
	if isInboxTag, ok := args["is_inbox_tag"].(bool); ok {
		tag.IsInboxTag = isInboxTag
	}

	// Call API
	createdTag, err := s.paperlessClient.CreateTag(ctx, tag)
	if err != nil {
		slog.Error("Failed to create tag", "name", name, "error", err)
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return createdTag, nil
}

// handleUpdateTag handles the update_tag tool
func (s *Server) handleUpdateTag(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract tag ID
	tagID, ok := args["tag_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("tag_id is required and must be an integer")
	}

	slog.Debug("Update tag tool invoked", "tag_id", int(tagID))

	// Build updates map
	updates := make(map[string]interface{})

	if name, ok := args["name"].(string); ok {
		updates["name"] = name
	}
	if color, ok := args["color"].(string); ok {
		updates["color"] = color
	}
	if match, ok := args["match"].(string); ok {
		updates["match"] = match
	}
	if matchingAlgo, ok := args["matching_algorithm"].(float64); ok {
		updates["matching_algorithm"] = int(matchingAlgo)
	}
	if isInsensitive, ok := args["is_insensitive"].(bool); ok {
		updates["is_insensitive"] = isInsensitive
	}
	if isInboxTag, ok := args["is_inbox_tag"].(bool); ok {
		updates["is_inbox_tag"] = isInboxTag
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field must be provided for update")
	}

	// Call API
	updatedTag, err := s.paperlessClient.UpdateTag(ctx, int(tagID), updates)
	if err != nil {
		slog.Error("Failed to update tag", "tag_id", int(tagID), "error", err)
		return nil, fmt.Errorf("failed to update tag: %w", err)
	}

	return updatedTag, nil
}

// handleDeleteTag handles the delete_tag tool
func (s *Server) handleDeleteTag(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract tag ID
	tagID, ok := args["tag_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("tag_id is required and must be an integer")
	}

	slog.Debug("Delete tag tool invoked", "tag_id", int(tagID))

	// Call API
	err := s.paperlessClient.DeleteTag(ctx, int(tagID))
	if err != nil {
		slog.Error("Failed to delete tag", "tag_id", int(tagID), "error", err)
		return nil, fmt.Errorf("failed to delete tag: %w", err)
	}

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Tag %d deleted successfully", int(tagID)),
	}, nil
}
