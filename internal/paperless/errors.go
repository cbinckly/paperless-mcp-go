package paperless

import (
	"fmt"
	"net/http"
)

// Error represents a Paperless API error
type Error struct {
	StatusCode int
	Message    string
	Details    map[string]interface{}
}

func (e *Error) Error() string {
	if len(e.Details) > 0 {
		return fmt.Sprintf("paperless API error (status %d): %s - %v",
			e.StatusCode, e.Message, e.Details)
	}
	return fmt.Sprintf("paperless API error (status %d): %s",
		e.StatusCode, e.Message)
}

// NewError creates a new API error
func NewError(statusCode int, message string, details map[string]interface{}) *Error {
	return &Error{
		StatusCode: statusCode,
		Message:    message,
		Details:    details,
	}
}

// IsNotFound checks if error is a 404
func IsNotFound(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.StatusCode == http.StatusNotFound
	}
	return false
}

// IsUnauthorized checks if error is a 401/403
func IsUnauthorized(err error) bool {
	if apiErr, ok := err.(*Error); ok {
		return apiErr.StatusCode == http.StatusUnauthorized ||
			apiErr.StatusCode == http.StatusForbidden
	}
	return false
}
