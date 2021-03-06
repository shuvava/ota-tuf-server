package db

import (
	"context"

	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// KeyRepository is the interface for the data.RepoKey repository.
type KeyRepository interface {
	// Create persist new data.Object in database
	Create(ctx context.Context, obj data.RepoKey) error
	// FindByRepoId returns data.RepoKey by repoId
	FindByRepoId(ctx context.Context, repoID data.RepoID) ([]data.RepoKey, error)
	// FindByKeyID returns data.RepoKey by keyID
	FindByKeyID(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (*data.RepoKey, error)
	// Exists checks if data.RepoKey exists in database
	Exists(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (bool, error)
}
