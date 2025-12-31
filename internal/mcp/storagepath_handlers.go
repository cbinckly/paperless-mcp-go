package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// handleListStoragePaths handles the list_storage_paths tool
func (s *Server) handleListStoragePaths(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	slog.Debug("Listing storage paths", "page", page, "page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.ListStoragePaths(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to list storage paths", "error", err)
		return nil, fmt.Errorf("failed to list storage paths: %w", err)
	}

	// Parse storage paths from Results
	var storagePaths []paperless.StoragePath
	if err := json.Unmarshal(response.Results, &storagePaths); err != nil {
		slog.Error("Failed to parse storage paths results", "error", err)
		return nil, fmt.Errorf("failed to parse results: %w", err)
	}

	slog.Info("Storage paths listed successfully",
		"count", response.Count,
		"returned", len(storagePaths))

	return map[string]interface{}{
		"count":         response.Count,
		"page":          page,
		"page_size":     pageSize,
		"has_next":      response.Next != nil,
		"has_prev":      response.Previous != nil,
		"storage_paths": storagePaths,
	}, nil
}

// handleGetStoragePath handles the get_storage_path tool
func (s *Server) handleGetStoragePath(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate storage_path_id
	storagePathIDFloat, ok := args["storage_path_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("storage_path_id parameter is required and must be an integer")
	}
	storagePathID := int(storagePathIDFloat)
	if storagePathID < 1 {
		return nil, fmt.Errorf("storage_path_id must be a positive integer")
	}

	slog.Debug("Getting storage path", "storage_path_id", storagePathID)

	// Call Paperless API
	storagePath, err := s.paperlessClient.GetStoragePath(ctx, storagePathID)
	if err != nil {
		slog.Error("Failed to get storage path",
			"storage_path_id", storagePathID,
			"error", err)
		return nil, fmt.Errorf("failed to get storage path: %w", err)
	}

	slog.Info("Storage path retrieved successfully",
		"storage_path_id", storagePathID,
		"name", storagePath.Name)

	return storagePath, nil
}

// handleCreateStoragePath handles the create_storage_path tool
func (s *Server) handleCreateStoragePath(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract required name
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required and must be a non-empty string")
	}

	// Extract required path
	pathStr, ok := args["path"].(string)
	if !ok || pathStr == "" {
		return nil, fmt.Errorf("path parameter is required and must be a non-empty string")
	}

	slog.Debug("Creating storage path", "name", name, "path", pathStr)

	// Build storage path from args
	storagePath := &paperless.StoragePath{
		Name: name,
		Path: pathStr,
	}

	// Extract optional fields
	if match, ok := args["match"].(string); ok {
		storagePath.Match = match
	}
	if matchingAlg, ok := args["matching_algorithm"].(float64); ok {
		storagePath.MatchingAlgorithm = int(matchingAlg)
	}
	if isInsensitive, ok := args["is_insensitive"].(bool); ok {
		storagePath.IsInsensitive = isInsensitive
	}

	// Call Paperless API
	createdStoragePath, err := s.paperlessClient.CreateStoragePath(ctx, storagePath)
	if err != nil {
		slog.Error("Failed to create storage path",
			"name", name,
			"error", err)
		return nil, fmt.Errorf("failed to create storage path: %w", err)
	}

	slog.Info("Storage path created successfully",
		"storage_path_id", createdStoragePath.ID,
		"name", createdStoragePath.Name)

	return createdStoragePath, nil
}

// handleUpdateStoragePath handles the update_storage_path tool
func (s *Server) handleUpdateStoragePath(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate storage_path_id
	storagePathIDFloat, ok := args["storage_path_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("storage_path_id parameter is required and must be an integer")
	}
	storagePathID := int(storagePathIDFloat)
	if storagePathID < 1 {
		return nil, fmt.Errorf("storage_path_id must be a positive integer")
	}

	// Build updates map from args (exclude storage_path_id)
	updates := make(map[string]interface{})
	for key, value := range args {
		if key != "storage_path_id" {
			updates[key] = value
		}
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field to update must be provided")
	}

	slog.Debug("Updating storage path",
		"storage_path_id", storagePathID,
		"fields", len(updates))

	// Call Paperless API
	updatedStoragePath, err := s.paperlessClient.UpdateStoragePath(ctx, storagePathID, updates)
	if err != nil {
		slog.Error("Failed to update storage path",
			"storage_path_id", storagePathID,
			"error", err)
		return nil, fmt.Errorf("failed to update storage path: %w", err)
	}

	slog.Info("Storage path updated successfully",
		"storage_path_id", storagePathID,
		"name", updatedStoragePath.Name)

	return updatedStoragePath, nil
}

// handleDeleteStoragePath handles the delete_storage_path tool
func (s *Server) handleDeleteStoragePath(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate storage_path_id
	storagePathIDFloat, ok := args["storage_path_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("storage_path_id parameter is required and must be an integer")
	}
	storagePathID := int(storagePathIDFloat)
	if storagePathID < 1 {
		return nil, fmt.Errorf("storage_path_id must be a positive integer")
	}

	slog.Debug("Deleting storage path", "storage_path_id", storagePathID)

	// Call Paperless API
	err := s.paperlessClient.DeleteStoragePath(ctx, storagePathID)
	if err != nil {
		slog.Error("Failed to delete storage path",
			"storage_path_id", storagePathID,
			"error", err)
		return nil, fmt.Errorf("failed to delete storage path: %w", err)
	}

	slog.Info("Storage path deleted successfully", "storage_path_id", storagePathID)

	return map[string]interface{}{
		"success":         true,
		"storage_path_id": storagePathID,
		"message":         "Storage path deleted successfully",
	}, nil
}
