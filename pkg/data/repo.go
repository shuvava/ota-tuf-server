package data

import (
	"github.com/shuvava/ota-tuf-server/pkg/encryption"

	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

// Repo is a TUF repo object
type Repo struct {
	Namespace      cmndata.Namespace  `json:"namespace"`
	RepoID         RepoID             `json:"repoId"`
	KeyType        encryption.KeyType `json:"keyType"`
	Threshold      uint               `json:"threshold"` // number of keys required to sign content to be valid
	CurrentVersion uint               `json:"currentVersion"`
}
