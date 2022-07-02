package data

import (
	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

// Repo is a TUF repo object
type Repo struct {
	Namespace cmndata.Namespace `json:"namespace"`
	// RepoID is the id of the repo
	RepoID RepoID `json:"repo_id"`
}
