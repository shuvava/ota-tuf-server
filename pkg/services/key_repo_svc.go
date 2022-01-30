package services

import (
	"context"

	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

type RepositoryService struct {
	log logger.Logger
	db  db.KeyRepository
}

// NewRepositoryService creates new instance of services.RepositoryService
func NewRepositoryService(l logger.Logger, db db.KeyRepository) *RepositoryService {
	log := l.SetOperation("repository-service")
	return &RepositoryService{
		log: log,
		db:  db,
	}
}

// CreateNewRepository initializes new repository by creating and persisting new key pair for data.TopLevelRoles
func (svc *RepositoryService) CreateNewRepository(ctx context.Context, repoID data.RepoID, keyType data.KeyType) error {
	keys := make([]data.RepoKey, len(data.TopLevelRoles))
	i := 0
	for role := range data.TopLevelRoles {
		key, err := encryption.NewKey(keyType)
		if err != nil {
			return err
		}
		keySerialized, err := key.MarshalAllData()
		keys[i] = data.RepoKey{
			RepoID: repoID,
			Role:   role,
			KeyID:  data.NewKeyID(repoID, role),
			Key:    *keySerialized,
		}
		i += 1
	}
	for _, key := range keys {
		err := svc.db.Create(ctx, key)
		if err != nil {
			return err
		}
	}
	return nil
}
