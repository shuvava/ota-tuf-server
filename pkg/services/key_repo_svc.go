package services

import (
	"context"

	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

// KeyRepositoryService is the service responsible for managing keys of repositories
type KeyRepositoryService struct {
	log logger.Logger
	db  db.KeyRepository
}

// NewKeyRepositoryService creates new instance of services.KeyRepositoryService
func NewKeyRepositoryService(l logger.Logger, db db.KeyRepository) *KeyRepositoryService {
	log := l.SetArea("key-repository-service")
	return &KeyRepositoryService{
		log: log,
		db:  db,
	}
}

// CreateNewKey creates new repository key
func (svc *KeyRepositoryService) CreateNewKey(ctx context.Context, repo *data.Repo, role data.RoleType) (*data.RepoKey, error) {
	key, err := encryption.NewKey(repo.KeyType)
	if err != nil {
		return nil, err
	}
	keySerialized, err := key.MarshalAllData()
	if err != nil {
		return nil, err
	}
	keyObj := &data.RepoKey{
		RepoID: repo.RepoID,
		Role:   role,
		KeyID:  data.NewKeyID(key),
		Key:    *keySerialized,
	}
	err = svc.db.Create(ctx, keyObj)

	return keyObj, err
}

// CreateAllRolesKeys create keys for all root roles
func (svc *KeyRepositoryService) CreateAllRolesKeys(ctx context.Context, repo *data.Repo) ([]*data.RepoKey, error) {
	keys := make([]*data.RepoKey, 0)
	for role := range data.TopLevelRoles {
		k, err := svc.CreateNewKey(ctx, repo, role)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// ExistsKeyRole checks if data.RepoKey with role exists in database
func (svc *KeyRepositoryService) ExistsKeyRole(ctx context.Context, repoID data.RepoID, role data.RoleType) (bool, error) {
	return svc.db.ExistsKeyRole(ctx, repoID, role)
}

// GetRepoKeys returns []data.RepoKey for provided repoID
func (svc *KeyRepositoryService) GetRepoKeys(ctx context.Context, repoID data.RepoID) ([]*data.RepoKey, error) {
	return svc.db.FindByRepoID(ctx, repoID)
}
