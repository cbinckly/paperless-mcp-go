package mcp

import (
	"context"
	"fmt"
	"log/slog"
)

// Tool execution error messages
const (
	ErrToolNotFound     = "tool not found: %s"
	ErrToolExecFailed   = "tool execution failed: %w"
)

// ExecuteTool executes a registered tool by name
func (s *Server) ExecuteTool(ctx context.Context, toolName string, args map[string]interface{}) (interface{}, error) {
	// Check if tool exists
	tool, exists := s.tools[toolName]
	if !exists {
		slog.Warn("Tool not found",
			"tool", toolName,
			"available_tools", s.getToolNames())
		return nil, fmt.Errorf(ErrToolNotFound, toolName)
	}

	// Log execution start
	slog.Debug("Executing tool",
		"tool", toolName,
		"args_count", len(args))

	// Execute the tool handler
	result, err := tool.Handler(ctx, args)
	if err != nil {
		slog.Error("Tool execution failed",
			"tool", toolName,
			"error", err)
		return nil, fmt.Errorf(ErrToolExecFailed, err)
	}

	// Log successful execution
	slog.Debug("Tool executed successfully",
		"tool", toolName)

	return result, nil
}

// getToolNames returns a list of all registered tool names
func (s *Server) getToolNames() []string {
	names := make([]string, 0, len(s.tools))
	for name := range s.tools {
		names = append(names, name)
	}
	return names
}

// GetToolCount returns the number of registered tools
func (s *Server) GetToolCount() int {
	return len(s.tools)
}

// HasTool checks if a tool is registered
func (s *Server) HasTool(name string) bool {
	_, exists := s.tools[name]
	return exists
}

// GetToolInfo returns information about a registered tool
func (s *Server) GetToolInfo(name string) (Tool, bool) {
	tool, exists := s.tools[name]
	return tool, exists
}

// ListTools returns all registered tools
func (s *Server) ListTools() []Tool {
	tools := make([]Tool, 0, len(s.tools))
	for _, tool := range s.tools {
		tools = append(tools, tool)
	}
	return tools
}
