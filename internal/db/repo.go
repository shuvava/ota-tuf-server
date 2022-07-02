package db

import (
	"context"

	cmndata "github.com/shuvava/go-ota-svc-common/data"

	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// TufRepoRepository is the interface for the data.RepoKey repository.
type TufRepoRepository interface {
	// Create persist new data.Object in database
	Create(ctx context.Context, obj data.Repo) error
	// FindByNamespace returns data.Repo by Namespace
	FindByNamespace(ctx context.Context, ns cmndata.Namespace) (data.Repo, error)
	// FindById returns data.Repo by RepoID
	FindById(ctx context.Context, id data.RepoID) (data.Repo, error)
	// Exists checks if data.Repo exists in database
	Exists(ctx context.Context, ns cmndata.Namespace) (bool, error)
	// List returns all data.Repo
	List(ctx context.Context, skip, limit int64, sortFiled *string) ([]*data.Repo, int64, error)
}
