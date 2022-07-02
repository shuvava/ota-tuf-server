package services

import (
	"context"

	"github.com/shuvava/go-logging/logger"
	cmndata "github.com/shuvava/go-ota-svc-common/data"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

// RepositoryService is the service responsible for managing the repository
type RepositoryService struct {
	log    logger.Logger
	db     db.TufRepoRepository
	keySvc *KeyRepositoryService
}

// NewRepositoryService creates new instance of services.RepositoryService
func NewRepositoryService(l logger.Logger, keySvc *KeyRepositoryService, db db.TufRepoRepository) *RepositoryService {
	log := l.SetArea("repository-service")
	return &RepositoryService{
		log:    log,
		db:     db,
		keySvc: keySvc,
	}
}

// Create initializes new repository by creating and persisting new key pair for data.TopLevelRoles
func (svc *RepositoryService) Create(ctx context.Context, ns cmndata.Namespace, repoID data.RepoID, keyType data.KeyType) error {
	//log := svc.log.SetOperation("Create").WithContext(ctx)
	repo := data.Repo{
		Namespace: ns,
		RepoID:    repoID,
	}
	if err := svc.db.Create(ctx, repo); err != nil {
		return err
	}
	for role := range data.TopLevelRoles {
		_, err := svc.keySvc.CreateNewKey(ctx, repoID, role, keyType)
		if err != nil {
			return err
		}
	}

	return nil
}

// List returns all data.Repo
func (svc *RepositoryService) List(ctx context.Context, skip, limit int64) ([]*data.Repo, int64, error) {
	return svc.db.List(ctx, skip, limit, nil)
}