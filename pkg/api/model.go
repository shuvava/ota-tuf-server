package api

import "github.com/shuvava/ota-tuf-server/pkg/data"

// PaginatedResponse is generic paginated response of TUF API
type PaginatedResponse[V data.Repo] struct {
	Data  []*V  `json:"data"`
	Total int64 `json:"total"`
}
