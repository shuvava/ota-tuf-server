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
func (svc *KeyRepositoryService) CreateNewKey(ctx context.Context, repoID data.RepoID, role data.RoleType, keyType data.KeyType) (data.KeyID, error) {
	key, err := encryption.NewKey(keyType)
	if err != nil {
		return "", err
	}
	keySerialized, err := key.MarshalAllData()
	if err != nil {
		return "", err
	}
	keyObj := data.RepoKey{
		RepoID: repoID,
		Role:   role,
		KeyID:  encryption.NewKeyID(key),
		Key:    *keySerialized,
	}
	err = svc.db.Create(ctx, keyObj)

	return keyObj.KeyID, err
}

// ExistsKeyRole checks if data.RepoKey with role exists in database
func (svc *KeyRepositoryService) ExistsKeyRole(ctx context.Context, repoID data.RepoID, role data.RoleType) (bool, error) {
	return svc.db.ExistsKeyRole(ctx, repoID, role)
}
