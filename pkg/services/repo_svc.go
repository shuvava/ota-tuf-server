package services

import (
	"context"
	"time"

	"github.com/shuvava/go-logging/logger"
	cmndata "github.com/shuvava/go-ota-svc-common/data"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

// RepositoryService is the service responsible for managing the repository
type RepositoryService struct {
	log          logger.Logger
	db           db.TufRepoRepository
	keySvc       *KeyRepositoryService
	signedCtnSvc *SignedContentService
}

// NewRepositoryService creates new instance of services.RepositoryService
func NewRepositoryService(l logger.Logger, keySvc *KeyRepositoryService, signedCtnSvc *SignedContentService, db db.TufRepoRepository) *RepositoryService {
	log := l.SetArea("repository-service")
	return &RepositoryService{
		log:          log,
		db:           db,
		keySvc:       keySvc,
		signedCtnSvc: signedCtnSvc,
	}
}

// Create initializes new repository by creating and persisting new key pair for data.TopLevelRoles
func (svc *RepositoryService) Create(ctx context.Context, ns cmndata.Namespace, repoID data.RepoID, keyType encryption.KeyType) ([]data.KeyID, error) {
	log := svc.log.SetOperation("Create").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("Namespace", ns).
		WithField("RepoID", repoID).
		Debug("Creating new TUF repository")
	repo := data.Repo{
		Namespace: ns,
		RepoID:    repoID,
	}
	if err := svc.db.Create(ctx, repo); err != nil {
		return nil, err
	}
	keys := make([]data.KeyID, 0)
	for role := range data.TopLevelRoles {
		k, err := svc.keySvc.CreateNewKey(ctx, repoID, role, keyType)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	if err := svc.signedCtnSvc.CreateNewRepoSignedMeta(ctx, repoID); err != nil {
		return nil, err
	}

	return keys, nil
}

// List returns all data.Repo
func (svc *RepositoryService) List(ctx context.Context, skip, limit int64) ([]*data.Repo, int64, error) {
	return svc.db.List(ctx, skip, limit, nil)
}

// FindByNamespace returns data.Repo assigned to the namespace or ErrorDbNoDocumentFound
func (svc *RepositoryService) FindByNamespace(ctx context.Context, ns cmndata.Namespace) (*data.Repo, error) {
	return svc.db.FindByNamespace(ctx, ns)
}
