package paperless

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"
)

// Date format constants
const (
	DateOnlyFormat = "2006-01-02"
)

// FlexibleTime is a time.Time wrapper that can parse multiple date/time formats
// It handles both RFC3339 timestamps and date-only strings from the Paperless API
type FlexibleTime struct {
	time.Time
}

// UnmarshalJSON implements custom JSON unmarshaling for flexible date parsing
func (ft *FlexibleTime) UnmarshalJSON(data []byte) error {
	// Remove quotes from JSON string
	str := strings.Trim(string(data), `"`)
	
	// Handle empty string or null
	if str == "" || str == "null" {
		ft.Time = time.Time{}
		return nil
	}

	// Try parsing as RFC3339 first (full timestamp with timezone)
	if t, err := time.Parse(time.RFC3339, str); err == nil {
		ft.Time = t
		slog.Debug("Parsed time as RFC3339", "input", str, "result", t)
		return nil
	}

	// Try parsing as date-only format
	if t, err := time.Parse(DateOnlyFormat, str); err == nil {
		ft.Time = t
		slog.Debug("Parsed time as date-only", "input", str, "result", t)
		return nil
	}

	// Both formats failed
	return fmt.Errorf("unable to parse time '%s' as RFC3339 or date-only format", str)
}

// MarshalJSON implements JSON marshaling, outputting RFC3339 format
func (ft FlexibleTime) MarshalJSON() ([]byte, error) {
	// Marshal as RFC3339 format string for consistency
	if ft.Time.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, ft.Time.Format(time.RFC3339))), nil
}

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
	Created             FlexibleTime         `json:"created"`
	CreatedDate         string               `json:"created_date"`
	Modified            FlexibleTime         `json:"modified"`
	Added               FlexibleTime         `json:"added"`
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
	ID                 int          `json:"id"`
	Slug               string       `json:"slug"`
	Name               string       `json:"name"`
	Match              string       `json:"match"`
	MatchingAlgorithm  int          `json:"matching_algorithm"`
	IsInsensitive      bool         `json:"is_insensitive"`
	DocumentCount      int          `json:"document_count"`
	LastCorrespondence FlexibleTime `json:"last_correspondence,omitempty"`
	Owner              int          `json:"owner,omitempty"`
	UserCanChange      bool         `json:"user_can_change,omitempty"`
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
	ID       int          `json:"id"`
	Note     string       `json:"note"`
	Created  FlexibleTime `json:"created"`
	Document int          `json:"document"`
	User     *int         `json:"user"`
}

// SearchResult represents a search result
type SearchResult struct {
	Documents []Document `json:"results"`
	Count     int        `json:"count"`
	Page      int        `json:"page,omitempty"`
	PageCount int        `json:"page_count,omitempty"`
}
