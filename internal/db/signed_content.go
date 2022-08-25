package db

import (
	"context"

	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// TufSignedContent is the interface for the data.SignedRootRole repository.
type TufSignedContent interface {
	// Exists checks if data.SignedRootRole exists in database
	Exists(ctx context.Context, repoID data.RepoID, ver uint) (bool, error)
	// Create persists new data.SignedRootRole object
	Create(ctx context.Context, role *data.SignedRootRole) error
	// FindVersion returns data.SignedRootRole with Version equal to ver parameter
	FindVersion(ctx context.Context, repoID data.RepoID, ver uint) (*data.SignedRootRole, error)
	// GetMaxVersion returns current repo version
	GetMaxVersion(ctx context.Context, repoID data.RepoID) (uint, error)
}
