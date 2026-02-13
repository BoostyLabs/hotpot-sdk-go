package client

// Page represents a paginated data response.
type Page[T any] struct {
	Data     []T          `json:"data"`
	Metadata PageMetadata `json:"pagination"`
}

// PageMetadata represents pagination metadata.
type PageMetadata struct {
	Total  int64 `json:"total"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
	Pages  int64 `json:"pages"`
}
