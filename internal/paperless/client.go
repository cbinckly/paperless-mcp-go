package paperless

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"
	"net/url"
)

// Client constants
const (
	DefaultTimeout    = 30 * time.Second
	AuthHeaderName    = "Authorization"
	AuthTokenPrefix   = "Token"
	ContentTypeJSON   = "application/json"
	ContentTypeHeader = "Content-Type"
)

// Pagination constants
const (
	DefaultPageSize = 25
	MaxPageSize     = 100
)


// Client represents a Paperless API client
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

// New creates a new Paperless API client
func New(baseURL, token string) *Client {
	return &Client{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		token:   token,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// doRequest performs an HTTP request with authentication
func (c *Client) doRequest(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	// Build full URL
	url := c.baseURL + path

	// Create request with context
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		slog.Error("Failed to create HTTP request",
			"method", method,
			"url", url,
			"error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authorization header
	req.Header.Set(AuthHeaderName, fmt.Sprintf("%s %s", AuthTokenPrefix, c.token))

	// Add content type for requests with body
	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		req.Header.Set(ContentTypeHeader, ContentTypeJSON)
	}

	// Log request (without sensitive data)
	slog.Debug("Making API request",
		"method", method,
		"url", url)

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		slog.Error("HTTP request failed",
			"method", method,
			"url", url,
			"error", err)
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Log response
	slog.Debug("Received API response",
		"method", method,
		"url", url,
		"status", resp.StatusCode)

	return resp, nil
}

// GET performs a GET request
func (c *Client) GET(ctx context.Context, path string) ([]byte, error) {
	resp, err := c.doRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body",
			"path", path,
			"error", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// POST performs a POST request
func (c *Client) POST(ctx context.Context, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			slog.Error("Failed to marshal request body",
				"path", path,
				"error", err)
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	resp, err := c.doRequest(ctx, http.MethodPost, path, bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body",
			"path", path,
			"error", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// PUT performs a PUT request
func (c *Client) PUT(ctx context.Context, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			slog.Error("Failed to marshal request body",
				"path", path,
				"error", err)
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	resp, err := c.doRequest(ctx, http.MethodPut, path, bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body",
			"path", path,
			"error", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// PATCH performs a PATCH request
func (c *Client) PATCH(ctx context.Context, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			slog.Error("Failed to marshal request body",
				"path", path,
				"error", err)
			return nil, fmt.Errorf("failed to marshal body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	resp, err := c.doRequest(ctx, http.MethodPatch, path, bodyReader)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read response body",
			"path", path,
			"error", err)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, parseError(resp.StatusCode, bodyBytes)
	}

	return bodyBytes, nil
}

// DELETE performs a DELETE request
func (c *Client) DELETE(ctx context.Context, path string) error {
	resp, err := c.doRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// DELETE can return 204 No Content or 200 OK
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return parseError(resp.StatusCode, bodyBytes)
	}

	return nil
}


// SearchDocuments searches for documents by text query with pagination
func (c *Client) SearchDocuments(ctx context.Context, query string, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Build query string
	path := fmt.Sprintf("/api/documents/?query=%s&page=%d&page_size=%d",
		url.QueryEscape(query), page, pageSize)

	slog.Debug("Searching documents",
		"query", query,
		"page", page,
		"page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse search response",
			"error", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetSimilarDocuments finds documents similar to a given document with pagination
func (c *Client) GetSimilarDocuments(ctx context.Context, documentID int, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	// Build path
	path := fmt.Sprintf("/api/documents/%d/similar/?page=%d&page_size=%d",
		documentID, page, pageSize)

	slog.Debug("Finding similar documents",
		"document_id", documentID,
		"page", page,
		"page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse similar documents response",
			"error", err)
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response, nil
}

// GetDocument retrieves a document by ID
func (c *Client) GetDocument(ctx context.Context, documentID int) (*Document, error) {
	path := fmt.Sprintf("/api/documents/%d/", documentID)

	slog.Debug("Getting document", "document_id", documentID)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var document Document
	if err := json.Unmarshal(bodyBytes, &document); err != nil {
		slog.Error("Failed to parse document response",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to parse document: %w", err)
	}

	return &document, nil
}

// GetDocumentContent retrieves the text content of a document
func (c *Client) GetDocumentContent(ctx context.Context, documentID int) (string, error) {
	// First get the document to access its content
	document, err := c.GetDocument(ctx, documentID)
	if err != nil {
		return "", err
	}

	slog.Debug("Retrieved document content",
		"document_id", documentID,
		"content_length", len(document.Content))

	return document.Content, nil
}

// CreateDocument creates a new document
func (c *Client) CreateDocument(ctx context.Context, document *Document) (*Document, error) {
	path := "/api/documents/"

	slog.Debug("Creating document", "title", document.Title)

	// Make POST request
	bodyBytes, err := c.POST(ctx, path, document)
	if err != nil {
		return nil, err
	}

	// Parse response
	var createdDocument Document
	if err := json.Unmarshal(bodyBytes, &createdDocument); err != nil {
		slog.Error("Failed to parse created document response",
			"error", err)
		return nil, fmt.Errorf("failed to parse created document: %w", err)
	}

	slog.Info("Document created successfully",
		"document_id", createdDocument.ID,
		"title", createdDocument.Title)

	return &createdDocument, nil
}

// UpdateDocument updates a document's metadata
func (c *Client) UpdateDocument(ctx context.Context, documentID int, updates map[string]interface{}) (*Document, error) {
	path := fmt.Sprintf("/api/documents/%d/", documentID)

	slog.Debug("Updating document",
		"document_id", documentID,
		"fields", len(updates))

	// Make PATCH request
	bodyBytes, err := c.PATCH(ctx, path, updates)
	if err != nil {
		return nil, err
	}

	// Parse response
	var updatedDocument Document
	if err := json.Unmarshal(bodyBytes, &updatedDocument); err != nil {
		slog.Error("Failed to parse updated document response",
			"document_id", documentID,
			"error", err)
		return nil, fmt.Errorf("failed to parse updated document: %w", err)
	}

	slog.Info("Document updated successfully",
		"document_id", documentID,
		"title", updatedDocument.Title)

	return &updatedDocument, nil
}

// DeleteDocument deletes a document by ID
func (c *Client) DeleteDocument(ctx context.Context, documentID int) error {
	path := fmt.Sprintf("/api/documents/%d/", documentID)

	slog.Debug("Deleting document", "document_id", documentID)

	// Make DELETE request
	err := c.DELETE(ctx, path)
	if err != nil {
		return err
	}

	slog.Info("Document deleted successfully", "document_id", documentID)
	return nil
}

// ListCorrespondents retrieves all correspondents with pagination
func (c *Client) ListCorrespondents(ctx context.Context, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	path := fmt.Sprintf("/api/correspondents/?page=%d&page_size=%d", page, pageSize)

	slog.Debug("Listing correspondents", "page", page, "page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse correspondents response", "error", err)
		return nil, fmt.Errorf("failed to parse correspondents: %w", err)
	}

	return &response, nil
}

// GetCorrespondent retrieves a correspondent by ID
func (c *Client) GetCorrespondent(ctx context.Context, correspondentID int) (*Correspondent, error) {
	path := fmt.Sprintf("/api/correspondents/%d/", correspondentID)

	slog.Debug("Getting correspondent", "correspondent_id", correspondentID)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var correspondent Correspondent
	if err := json.Unmarshal(bodyBytes, &correspondent); err != nil {
		slog.Error("Failed to parse correspondent response",
			"correspondent_id", correspondentID,
			"error", err)
		return nil, fmt.Errorf("failed to parse correspondent: %w", err)
	}

	return &correspondent, nil
}

// CreateCorrespondent creates a new correspondent
func (c *Client) CreateCorrespondent(ctx context.Context, correspondent *Correspondent) (*Correspondent, error) {
	path := "/api/correspondents/"

	slog.Debug("Creating correspondent", "name", correspondent.Name)

	// Make POST request
	bodyBytes, err := c.POST(ctx, path, correspondent)
	if err != nil {
		return nil, err
	}

	// Parse response
	var createdCorrespondent Correspondent
	if err := json.Unmarshal(bodyBytes, &createdCorrespondent); err != nil {
		slog.Error("Failed to parse created correspondent response", "error", err)
		return nil, fmt.Errorf("failed to parse created correspondent: %w", err)
	}

	slog.Info("Correspondent created successfully",
		"correspondent_id", createdCorrespondent.ID,
		"name", createdCorrespondent.Name)

	return &createdCorrespondent, nil
}

// UpdateCorrespondent updates a correspondent's information
func (c *Client) UpdateCorrespondent(ctx context.Context, correspondentID int, updates map[string]interface{}) (*Correspondent, error) {
	path := fmt.Sprintf("/api/correspondents/%d/", correspondentID)

	slog.Debug("Updating correspondent",
		"correspondent_id", correspondentID,
		"fields", len(updates))

	// Make PATCH request
	bodyBytes, err := c.PATCH(ctx, path, updates)
	if err != nil {
		return nil, err
	}

	// Parse response
	var updatedCorrespondent Correspondent
	if err := json.Unmarshal(bodyBytes, &updatedCorrespondent); err != nil {
		slog.Error("Failed to parse updated correspondent response",
			"correspondent_id", correspondentID,
			"error", err)
		return nil, fmt.Errorf("failed to parse updated correspondent: %w", err)
	}

	slog.Info("Correspondent updated successfully",
		"correspondent_id", correspondentID,
		"name", updatedCorrespondent.Name)

	return &updatedCorrespondent, nil
}

// DeleteCorrespondent deletes a correspondent by ID
func (c *Client) DeleteCorrespondent(ctx context.Context, correspondentID int) error {
	path := fmt.Sprintf("/api/correspondents/%d/", correspondentID)

	slog.Debug("Deleting correspondent", "correspondent_id", correspondentID)

	// Make DELETE request
	err := c.DELETE(ctx, path)
	if err != nil {
		return err
	}

	slog.Info("Correspondent deleted successfully", "correspondent_id", correspondentID)
	return nil
}


// ListDocumentTypes retrieves all document types with pagination
func (c *Client) ListDocumentTypes(ctx context.Context, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	path := fmt.Sprintf("/api/document_types/?page=%d&page_size=%d", page, pageSize)

	slog.Debug("Listing document types", "page", page, "page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse document types response", "error", err)
		return nil, fmt.Errorf("failed to parse document types: %w", err)
	}

	return &response, nil
}

// GetDocumentType retrieves a document type by ID
func (c *Client) GetDocumentType(ctx context.Context, typeID int) (*DocumentType, error) {
	path := fmt.Sprintf("/api/document_types/%d/", typeID)

	slog.Debug("Getting document type", "document_type_id", typeID)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var docType DocumentType
	if err := json.Unmarshal(bodyBytes, &docType); err != nil {
		slog.Error("Failed to parse document type response",
			"document_type_id", typeID,
			"error", err)
		return nil, fmt.Errorf("failed to parse document type: %w", err)
	}

	return &docType, nil
}

// CreateDocumentType creates a new document type
func (c *Client) CreateDocumentType(ctx context.Context, docType *DocumentType) (*DocumentType, error) {
	path := "/api/document_types/"

	slog.Debug("Creating document type", "name", docType.Name)

	// Make POST request
	bodyBytes, err := c.POST(ctx, path, docType)
	if err != nil {
		return nil, err
	}

	// Parse response
	var createdDocType DocumentType
	if err := json.Unmarshal(bodyBytes, &createdDocType); err != nil {
		slog.Error("Failed to parse created document type response", "error", err)
		return nil, fmt.Errorf("failed to parse created document type: %w", err)
	}

	slog.Info("Document type created successfully",
		"document_type_id", createdDocType.ID,
		"name", createdDocType.Name)

	return &createdDocType, nil
}

// UpdateDocumentType updates a document type's information
func (c *Client) UpdateDocumentType(ctx context.Context, typeID int, updates map[string]interface{}) (*DocumentType, error) {
	path := fmt.Sprintf("/api/document_types/%d/", typeID)

	slog.Debug("Updating document type",
		"document_type_id", typeID,
		"fields", len(updates))

	// Make PATCH request
	bodyBytes, err := c.PATCH(ctx, path, updates)
	if err != nil {
		return nil, err
	}

	// Parse response
	var updatedDocType DocumentType
	if err := json.Unmarshal(bodyBytes, &updatedDocType); err != nil {
		slog.Error("Failed to parse updated document type response",
			"document_type_id", typeID,
			"error", err)
		return nil, fmt.Errorf("failed to parse updated document type: %w", err)
	}

	slog.Info("Document type updated successfully",
		"document_type_id", typeID,
		"name", updatedDocType.Name)

	return &updatedDocType, nil
}

// DeleteDocumentType deletes a document type by ID
func (c *Client) DeleteDocumentType(ctx context.Context, typeID int) error {
	path := fmt.Sprintf("/api/document_types/%d/", typeID)

	slog.Debug("Deleting document type", "document_type_id", typeID)

	// Make DELETE request
	err := c.DELETE(ctx, path)
	if err != nil {
		return err
	}

	slog.Info("Document type deleted successfully", "document_type_id", typeID)
	return nil
}


// ListTags retrieves all tags with pagination
func (c *Client) ListTags(ctx context.Context, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	path := fmt.Sprintf("/api/tags/?page=%d&page_size=%d", page, pageSize)

	slog.Debug("Listing tags", "page", page, "page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse tags response", "error", err)
		return nil, fmt.Errorf("failed to parse tags: %w", err)
	}

	return &response, nil
}

// GetTag retrieves a tag by ID
func (c *Client) GetTag(ctx context.Context, tagID int) (*Tag, error) {
	path := fmt.Sprintf("/api/tags/%d/", tagID)

	slog.Debug("Getting tag", "tag_id", tagID)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var tag Tag
	if err := json.Unmarshal(bodyBytes, &tag); err != nil {
		slog.Error("Failed to parse tag response",
			"tag_id", tagID,
			"error", err)
		return nil, fmt.Errorf("failed to parse tag: %w", err)
	}

	return &tag, nil
}

// CreateTag creates a new tag
func (c *Client) CreateTag(ctx context.Context, tag *Tag) (*Tag, error) {
	path := "/api/tags/"

	slog.Debug("Creating tag", "name", tag.Name)

	// Make POST request
	bodyBytes, err := c.POST(ctx, path, tag)
	if err != nil {
		return nil, err
	}

	// Parse response
	var createdTag Tag
	if err := json.Unmarshal(bodyBytes, &createdTag); err != nil {
		slog.Error("Failed to parse created tag response", "error", err)
		return nil, fmt.Errorf("failed to parse created tag: %w", err)
	}

	slog.Info("Tag created successfully",
		"tag_id", createdTag.ID,
		"name", createdTag.Name)

	return &createdTag, nil
}

// UpdateTag updates a tag's information
func (c *Client) UpdateTag(ctx context.Context, tagID int, updates map[string]interface{}) (*Tag, error) {
	path := fmt.Sprintf("/api/tags/%d/", tagID)

	slog.Debug("Updating tag",
		"tag_id", tagID,
		"fields", len(updates))

	// Make PATCH request
	bodyBytes, err := c.PATCH(ctx, path, updates)
	if err != nil {
		return nil, err
	}

	// Parse response
	var updatedTag Tag
	if err := json.Unmarshal(bodyBytes, &updatedTag); err != nil {
		slog.Error("Failed to parse updated tag response",
			"tag_id", tagID,
			"error", err)
		return nil, fmt.Errorf("failed to parse updated tag: %w", err)
	}

	slog.Info("Tag updated successfully",
		"tag_id", tagID,
		"name", updatedTag.Name)

	return &updatedTag, nil
}

// DeleteTag deletes a tag by ID
func (c *Client) DeleteTag(ctx context.Context, tagID int) error {
	path := fmt.Sprintf("/api/tags/%d/", tagID)

	slog.Debug("Deleting tag", "tag_id", tagID)

	// Make DELETE request
	err := c.DELETE(ctx, path)
	if err != nil {
		return err
	}

	slog.Info("Tag deleted successfully", "tag_id", tagID)
	return nil
}


// parseError parses an error response from the API
func parseError(statusCode int, body []byte) error {
	var errorData map[string]interface{}
	if err := json.Unmarshal(body, &errorData); err != nil {
		// If we can't parse as JSON, use the raw body as message
		return NewError(statusCode, string(body), nil)
	}

	// Try to extract common error message fields
	var message string
	if msg, ok := errorData["detail"].(string); ok {
		message = msg
	} else if msg, ok := errorData["message"].(string); ok {
		message = msg
	} else if msg, ok := errorData["error"].(string); ok {
		message = msg
	} else {
		message = "API request failed"
	}

	return NewError(statusCode, message, errorData)
}


// ListStoragePaths retrieves all storage paths with pagination
func (c *Client) ListStoragePaths(ctx context.Context, page, pageSize int) (*PaginatedResponse, error) {
	// Validate and set defaults for pagination
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = DefaultPageSize
	} else if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	path := fmt.Sprintf("/api/storage_paths/?page=%d&page_size=%d", page, pageSize)

	slog.Debug("Listing storage paths", "page", page, "page_size", pageSize)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var response PaginatedResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		slog.Error("Failed to parse storage paths response", "error", err)
		return nil, fmt.Errorf("failed to parse storage paths: %w", err)
	}

	return &response, nil
}

// GetStoragePath retrieves a storage path by ID
func (c *Client) GetStoragePath(ctx context.Context, pathID int) (*StoragePath, error) {
	path := fmt.Sprintf("/api/storage_paths/%d/", pathID)

	slog.Debug("Getting storage path", "path_id", pathID)

	// Make GET request
	bodyBytes, err := c.GET(ctx, path)
	if err != nil {
		return nil, err
	}

	// Parse response
	var storagePath StoragePath
	if err := json.Unmarshal(bodyBytes, &storagePath); err != nil {
		slog.Error("Failed to parse storage path response",
			"path_id", pathID,
			"error", err)
		return nil, fmt.Errorf("failed to parse storage path: %w", err)
	}

	return &storagePath, nil
}

// CreateStoragePath creates a new storage path
func (c *Client) CreateStoragePath(ctx context.Context, storagePath *StoragePath) (*StoragePath, error) {
	path := "/api/storage_paths/"

	slog.Debug("Creating storage path", "name", storagePath.Name)

	// Make POST request
	bodyBytes, err := c.POST(ctx, path, storagePath)
	if err != nil {
		return nil, err
	}

	// Parse response
	var createdStoragePath StoragePath
	if err := json.Unmarshal(bodyBytes, &createdStoragePath); err != nil {
		slog.Error("Failed to parse created storage path response", "error", err)
		return nil, fmt.Errorf("failed to parse created storage path: %w", err)
	}

	slog.Info("Storage path created successfully",
		"path_id", createdStoragePath.ID,
		"name", createdStoragePath.Name)

	return &createdStoragePath, nil
}

// UpdateStoragePath updates a storage path's information
func (c *Client) UpdateStoragePath(ctx context.Context, pathID int, updates map[string]interface{}) (*StoragePath, error) {
	path := fmt.Sprintf("/api/storage_paths/%d/", pathID)

	slog.Debug("Updating storage path",
		"path_id", pathID,
		"fields", len(updates))

	// Make PATCH request
	bodyBytes, err := c.PATCH(ctx, path, updates)
	if err != nil {
		return nil, err
	}

	// Parse response
	var updatedStoragePath StoragePath
	if err := json.Unmarshal(bodyBytes, &updatedStoragePath); err != nil {
		slog.Error("Failed to parse updated storage path response",
			"path_id", pathID,
			"error", err)
		return nil, fmt.Errorf("failed to parse updated storage path: %w", err)
	}

	slog.Info("Storage path updated successfully",
		"path_id", pathID,
		"name", updatedStoragePath.Name)

	return &updatedStoragePath, nil
}

// DeleteStoragePath deletes a storage path by ID
func (c *Client) DeleteStoragePath(ctx context.Context, pathID int) error {
	path := fmt.Sprintf("/api/storage_paths/%d/", pathID)

	slog.Debug("Deleting storage path", "path_id", pathID)

	// Make DELETE request
	err := c.DELETE(ctx, path)
	if err != nil {
		return err
	}

	slog.Info("Storage path deleted successfully", "path_id", pathID)
	return nil
}
