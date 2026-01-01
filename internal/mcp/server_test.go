package mcp

import (
	"context"
	"encoding/json"
	"testing"

	"git.binckly.ca/cbinckly/paperless-mcp-go/internal/config"
)

// TestToolRegistrationWithSchema tests that tools can be registered with input schemas
// without causing the "both InputSchema and RawInputSchema set" error
func TestToolRegistrationWithSchema(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		PaperlessURL:   "http://localhost:8000",
		PaperlessToken: "test-token",
		MCPTransport:   "stdio",
	}

	// Create a new server - this will register all tools
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Verify that tools were registered
	if len(server.tools) == 0 {
		t.Fatal("No tools were registered")
	}

	t.Logf("Successfully registered %d tools", len(server.tools))

	// Verify specific tools have schemas
	expectedToolsWithSchemas := []string{
		"ping",
		"server_info",
		"search_documents",
		"get_document",
		"bulk_edit_documents",
	}

	for _, toolName := range expectedToolsWithSchemas {
		tool, exists := server.tools[toolName]
		if !exists {
			t.Errorf("Expected tool %s to be registered", toolName)
			continue
		}
		if tool.InputSchema == nil {
			t.Errorf("Expected tool %s to have an InputSchema", toolName)
		}
		t.Logf("Tool %s has InputSchema: %v", toolName, tool.InputSchema != nil)
	}
}

// TestToolSchemaMarshaling tests that tool schemas can be marshaled to JSON correctly
func TestToolSchemaMarshaling(t *testing.T) {
	// Test schema similar to what we use in tools.go
	schema := map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"query": map[string]interface{}{
				"type":        "string",
				"description": "Search query text",
			},
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number",
			},
		},
		"required": []string{"query"},
	}

	// Marshal to JSON
	schemaJSON, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	t.Logf("Schema JSON: %s", string(schemaJSON))

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(schemaJSON, &parsed); err != nil {
		t.Fatalf("Failed to parse schema JSON: %v", err)
	}

	// Verify structure
	if parsed["type"] != "object" {
		t.Errorf("Expected type to be 'object', got %v", parsed["type"])
	}
}

// TestToolHandlerExecution tests that tool handlers can be executed
func TestToolHandlerExecution(t *testing.T) {
	// Create a test config
	cfg := &config.Config{
		PaperlessURL:   "http://localhost:8000",
		PaperlessToken: "test-token",
		MCPTransport:   "stdio",
	}

	// Create a new server
	server, err := New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Test the ping tool directly
	ctx := context.Background()
	result, err := server.ExecuteTool(ctx, "ping", map[string]interface{}{})
	if err != nil {
		t.Fatalf("Failed to execute ping tool: %v", err)
	}

	// Verify result
	resultMap, ok := result.(map[string]string)
	if !ok {
		t.Fatalf("Expected result to be map[string]string, got %T", result)
	}

	if resultMap["status"] != "ok" {
		t.Errorf("Expected status 'ok', got %s", resultMap["status"])
	}

	if resultMap["message"] != "pong" {
		t.Errorf("Expected message 'pong', got %s", resultMap["message"])
	}

	t.Logf("Ping result: %+v", resultMap)
}
