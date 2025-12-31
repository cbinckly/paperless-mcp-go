package paperless

import (
	"encoding/json"
	"time"
)

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Count    int             `json:"count"`
	Next     *string         `json:"next"`
	Previous *string         `json:"previous"`
	All      []int           `json:"all,omitempty"`
	Results  json.RawMessage `json:"results"`
}

// Document represents a document in Paperless
type Document struct {
	ID                  int                  `json:"id"`
	Correspondent       *int                 `json:"correspondent"`
	DocumentType        *int                 `json:"document_type"`
	StoragePath         *int                 `json:"storage_path"`
	Title               string               `json:"title"`
	Content             string               `json:"content,omitempty"`
	Tags                []int                `json:"tags"`
	Created             time.Time            `json:"created"`
	CreatedDate         string               `json:"created_date"`
	Modified            time.Time            `json:"modified"`
	Added               time.Time            `json:"added"`
	ArchiveSerialNumber *int                 `json:"archive_serial_number"`
	OriginalFileName    string               `json:"original_file_name"`
	ArchivedFileName    *string              `json:"archived_file_name"`
	Owner               int                  `json:"owner,omitempty"`
	UserCanChange       bool                 `json:"user_can_change,omitempty"`
	Notes               []Note               `json:"notes,omitempty"`
	CustomFields        []CustomFieldValue   `json:"custom_fields,omitempty"`
}

// Correspondent represents a document correspondent
type Correspondent struct {
	ID                 int       `json:"id"`
	Slug               string    `json:"slug"`
	Name               string    `json:"name"`
	Match              string    `json:"match"`
	MatchingAlgorithm  int       `json:"matching_algorithm"`
	IsInsensitive      bool      `json:"is_insensitive"`
	DocumentCount      int       `json:"document_count"`
	LastCorrespondence time.Time `json:"last_correspondence,omitempty"`
	Owner              int       `json:"owner,omitempty"`
	UserCanChange      bool      `json:"user_can_change,omitempty"`
}

// DocumentType represents a document type
type DocumentType struct {
	ID                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Match             string `json:"match"`
	MatchingAlgorithm int    `json:"matching_algorithm"`
	IsInsensitive     bool   `json:"is_insensitive"`
	DocumentCount     int    `json:"document_count"`
	Owner             int    `json:"owner,omitempty"`
	UserCanChange     bool   `json:"user_can_change,omitempty"`
}

// Tag represents a document tag
type Tag struct {
	ID                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Color             string `json:"color"`
	Match             string `json:"match"`
	MatchingAlgorithm int    `json:"matching_algorithm"`
	IsInsensitive     bool   `json:"is_insensitive"`
	IsInboxTag        bool   `json:"is_inbox_tag"`
	DocumentCount     int    `json:"document_count"`
	Owner             int    `json:"owner,omitempty"`
	UserCanChange     bool   `json:"user_can_change,omitempty"`
}

// StoragePath represents a storage path
type StoragePath struct {
	ID                int    `json:"id"`
	Slug              string `json:"slug"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	Match             string `json:"match"`
	MatchingAlgorithm int    `json:"matching_algorithm"`
	IsInsensitive     bool   `json:"is_insensitive"`
	DocumentCount     int    `json:"document_count"`
	Owner             int    `json:"owner,omitempty"`
	UserCanChange     bool   `json:"user_can_change,omitempty"`
}

// CustomField represents a custom field definition
type CustomField struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	DataType string `json:"data_type"`
}

// CustomFieldValue represents a custom field value on a document
type CustomFieldValue struct {
	Field int         `json:"field"`
	Value interface{} `json:"value"`
}

// Note represents a document note
type Note struct {
	ID       int       `json:"id"`
	Note     string    `json:"note"`
	Created  time.Time `json:"created"`
	Document int       `json:"document"`
	User     *int      `json:"user"`
}

// SearchResult represents a search result
type SearchResult struct {
	Documents []Document `json:"results"`
	Count     int        `json:"count"`
	Page      int        `json:"page,omitempty"`
	PageCount int        `json:"page_count,omitempty"`
}
