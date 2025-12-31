package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/paperless"
)

// handleListCustomFields handles the list_custom_fields tool
func (s *Server) handleListCustomFields(ctx context.Context, args map[string]interface{}) (interface{}, error) {
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

	slog.Debug("Listing custom fields",
		"page", page,
		"page_size", pageSize)

	// Call Paperless API
	response, err := s.paperlessClient.ListCustomFields(ctx, page, pageSize)
	if err != nil {
		slog.Error("Failed to list custom fields", "error", err)
		return nil, fmt.Errorf("failed to list custom fields: %w", err)
	}

	// Parse custom fields from Results
	var fields []paperless.CustomField
	if err := json.Unmarshal(response.Results, &fields); err != nil {
		slog.Error("Failed to parse custom fields list", "error", err)
		return nil, fmt.Errorf("failed to parse custom fields: %w", err)
	}

	slog.Info("Custom fields listed successfully",
		"count", response.Count,
		"returned", len(fields))

	return map[string]interface{}{
		"count":      response.Count,
		"page":       page,
		"page_size":  pageSize,
		"has_next":   response.Next != nil,
		"has_prev":   response.Previous != nil,
		"fields":     fields,
	}, nil
}

// handleGetCustomField handles the get_custom_field tool
func (s *Server) handleGetCustomField(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate field_id
	fieldIDFloat, ok := args["field_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("field_id parameter is required and must be an integer")
	}
	fieldID := int(fieldIDFloat)
	if fieldID < 1 {
		return nil, fmt.Errorf("field_id must be a positive integer")
	}

	slog.Debug("Getting custom field", "field_id", fieldID)

	// Call Paperless API
	field, err := s.paperlessClient.GetCustomField(ctx, fieldID)
	if err != nil {
		slog.Error("Failed to get custom field",
			"field_id", fieldID,
			"error", err)
		return nil, fmt.Errorf("failed to get custom field: %w", err)
	}

	slog.Info("Custom field retrieved successfully",
		"field_id", fieldID,
		"name", field.Name)

	return field, nil
}

// handleCreateCustomField handles the create_custom_field tool
func (s *Server) handleCreateCustomField(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate required fields
	name, ok := args["name"].(string)
	if !ok || name == "" {
		return nil, fmt.Errorf("name parameter is required and must be a non-empty string")
	}

	dataType, ok := args["data_type"].(string)
	if !ok || dataType == "" {
		return nil, fmt.Errorf("data_type parameter is required and must be a non-empty string")
	}

	slog.Debug("Creating custom field",
		"name", name,
		"data_type", dataType)

	// Create custom field object
	field := &paperless.CustomField{
		Name:     name,
		DataType: dataType,
	}

	// Call Paperless API
	createdField, err := s.paperlessClient.CreateCustomField(ctx, field)
	if err != nil {
		slog.Error("Failed to create custom field",
			"name", name,
			"error", err)
		return nil, fmt.Errorf("failed to create custom field: %w", err)
	}

	slog.Info("Custom field created successfully",
		"field_id", createdField.ID,
		"name", createdField.Name)

	return createdField, nil
}

// handleUpdateCustomField handles the update_custom_field tool
func (s *Server) handleUpdateCustomField(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate field_id
	fieldIDFloat, ok := args["field_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("field_id parameter is required and must be an integer")
	}
	fieldID := int(fieldIDFloat)
	if fieldID < 1 {
		return nil, fmt.Errorf("field_id must be a positive integer")
	}

	// Build updates map with optional fields
	updates := make(map[string]interface{})
	
	if name, ok := args["name"].(string); ok && name != "" {
		updates["name"] = name
	}
	
	if dataType, ok := args["data_type"].(string); ok && dataType != "" {
		updates["data_type"] = dataType
	}

	if len(updates) == 0 {
		return nil, fmt.Errorf("at least one field to update must be provided")
	}

	slog.Debug("Updating custom field",
		"field_id", fieldID,
		"fields", len(updates))

	// Call Paperless API
	field, err := s.paperlessClient.UpdateCustomField(ctx, fieldID, updates)
	if err != nil {
		slog.Error("Failed to update custom field",
			"field_id", fieldID,
			"error", err)
		return nil, fmt.Errorf("failed to update custom field: %w", err)
	}

	slog.Info("Custom field updated successfully",
		"field_id", fieldID,
		"name", field.Name)

	return field, nil
}

// handleDeleteCustomField handles the delete_custom_field tool
func (s *Server) handleDeleteCustomField(ctx context.Context, args map[string]interface{}) (interface{}, error) {
	// Extract and validate field_id
	fieldIDFloat, ok := args["field_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("field_id parameter is required and must be an integer")
	}
	fieldID := int(fieldIDFloat)
	if fieldID < 1 {
		return nil, fmt.Errorf("field_id must be a positive integer")
	}

	slog.Debug("Deleting custom field", "field_id", fieldID)

	// Call Paperless API
	err := s.paperlessClient.DeleteCustomField(ctx, fieldID)
	if err != nil {
		slog.Error("Failed to delete custom field",
			"field_id", fieldID,
			"error", err)
		return nil, fmt.Errorf("failed to delete custom field: %w", err)
	}

	slog.Info("Custom field deleted successfully", "field_id", fieldID)

	return map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Custom field %d deleted successfully", fieldID),
	}, nil
}
