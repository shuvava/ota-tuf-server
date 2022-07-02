package data

import (
	"encoding/json"
	"time"
)

// SignedRootRole is a TUF repo signed Role
type SignedRootRole struct {
	// RepoID is the id of the repo
	RepoID    RepoID          `json:"repo_id"`
	ExpiresAt time.Time       `json:"expires_at"`
	Version   int             `json:"version"`
	Content   json.RawMessage `json:"signed_payload"`
}
