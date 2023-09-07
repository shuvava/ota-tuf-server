package services

import (
	"context"
	"errors"
	"time"

	intData "github.com/shuvava/ota-tuf-server/internal/data"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	cmndata "github.com/shuvava/go-ota-svc-common/data"
)

// RepositoryService is the service responsible for managing the repository
type RepositoryService struct {
	log          logger.Logger
	db           db.TufRepoRepository
	keySvc       *KeyRepositoryService
	signedCtnSvc *RepoVersionService
}

// NewRepositoryService creates new instance of services.RepositoryService
func NewRepositoryService(l logger.Logger, keySvc *KeyRepositoryService, signedCtnSvc *RepoVersionService, db db.TufRepoRepository) *RepositoryService {
	log := l.SetArea("repository-service")
	return &RepositoryService{
		log:          log,
		db:           db,
		keySvc:       keySvc,
		signedCtnSvc: signedCtnSvc,
	}
}

// Create initializes new repository by creating and persisting new key pair for data.TopLevelRoles
func (svc *RepositoryService) Create(ctx context.Context, ns cmndata.Namespace, repoID data.RepoID, keyType encryption.KeyType, threshold uint) ([]data.KeyID, error) {
	log := svc.log.SetOperation("Create").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField(intData.LogFieldNamespace, ns).
		WithField(intData.LogFieldRepoID, repoID).
		Debug("Creating new TUF repository")
	repo := &data.Repo{
		Namespace: ns,
		RepoID:    repoID,
		KeyType:   keyType,
		Threshold: threshold,
	}
	if err := svc.db.Create(ctx, repo); err != nil {
		return nil, err
	}
	keys, err := svc.keySvc.CreateAllRolesKeys(ctx, repo)
	if err != nil {
		return nil, err
	}

	if err = svc.signedCtnSvc.CreateFirstVersion(ctx, repoID, keys, threshold); err != nil {
		return nil, err
	}
	keyIds := make([]data.KeyID, len(keys))
	for i, k := range keys {
		keyIds[i] = k.KeyID
	}

	return keyIds, nil
}

// List returns all data.Repo
func (svc *RepositoryService) List(ctx context.Context, skip, limit int64) ([]*data.Repo, int64, error) {
	return svc.db.List(ctx, skip, limit, nil)
}

// FindByNamespace returns data.Repo assigned to the namespace or ErrorDbNoDocumentFound
func (svc *RepositoryService) FindByNamespace(ctx context.Context, ns cmndata.Namespace) (*data.Repo, error) {
	return svc.db.FindByNamespace(ctx, ns)
}

// Exists checks if TUF repository with data.RepoID exists
func (svc *RepositoryService) Exists(ctx context.Context, repoID data.RepoID) (bool, error) {
	_, err := svc.db.FindByID(ctx, repoID)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetAndRefresh returns repo updated signature of the TUF repository with data.RepoID
func (svc *RepositoryService) GetAndRefresh(ctx context.Context, repoID data.RepoID) (*data.SignedPayload[data.RepoSigned], error) {
	log := svc.log.SetOperation("GetAndRefresh").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField(intData.LogFieldRepoID, repoID).
		Debug("Get and refresh TUF repo")
	sig, err := svc.signedCtnSvc.GetCurrentVersion(ctx, repoID)
	if err != nil {
		return nil, err
	}
	if sig.ShouldBeRenewed() {
		repo, err := svc.db.FindByID(ctx, repoID)
		if err != nil {
			return nil, err
		}
		keys, err := svc.keySvc.GetRepoKeys(ctx, repoID)
		if err != nil {
			return nil, err
		}
		// generate new keys
		newKeys, err := svc.keySvc.CreateAllRolesKeys(ctx, repo)
		if err != nil {
			return nil, err
		}
		keys = append(keys, newKeys...)
		// and update repo signature
		sig, err = svc.signedCtnSvc.UpdateRepoVersion(ctx, sig, keys)
		if err != nil {
			return nil, err
		}
	}

	return sig.Content, nil
}

// SignPayload sign payload by provided role
func (svc *RepositoryService) SignPayload(ctx context.Context, repoID data.RepoID, role data.RoleType, payload interface{}) (*data.SignedPayload[data.RoleSign], error) {
	log := svc.log.SetOperation("SignPayload").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID).
		WithField(intData.LogFieldRole, role)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Signing payload with role keys")
	repo, err := svc.db.FindByID(ctx, repoID)
	if err != nil {
		return nil, err
	}
	keys, err := svc.keySvc.GetRepoKeysForRole(ctx, repoID, role)
	if err != nil {
		return nil, err
	}
	sig, err := data.SignPayload(payload, keys)
	if err != nil {
		return nil, err
	}
	roleSig := &data.RoleSign{
		Role:      role,
		Threshold: repo.Threshold,
	}
	sigPayload := &data.SignedPayload[data.RoleSign]{
		Signatures: sig,
		Signed:     roleSig,
	}
	return sigPayload, nil
}
