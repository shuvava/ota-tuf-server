package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

const (
	defaultRoleExpire = time.Hour * 24 * 365
)

// SignedContentService controls TUF SerializedKey roles life cycle
type SignedContentService struct {
	log    logger.Logger
	db     db.TufSignedContent
	keySvc *KeyRepositoryService
}

// NewSignedContentService creates new instance of services.SignedContentService
func NewSignedContentService(l logger.Logger, keySvc *KeyRepositoryService, db db.TufSignedContent) *SignedContentService {
	log := l.SetArea("signed-content-service")
	return &SignedContentService{
		log:    log,
		db:     db,
		keySvc: keySvc,
	}
}

func (svc *SignedContentService) GetRepoSignedMeta(ctx context.Context, repoID data.RepoID) (*data.SignedPayload[data.RootRole], error) {
	sig, err := svc.getCurrent(ctx, repoID)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == ErrorMissingSignedRole {
			return svc.createAndPersist(ctx, repoID)
		} else {
			return nil, err
		}
	}

	// TODO: implement
	return nil, fmt.Errorf("not Implemented")
}

// getCurrent returns the latest created data.SignedRootRole for data.Repo
func (svc *SignedContentService) getCurrent(ctx context.Context, repoID data.RepoID) (*data.SignedRootRole, error) {
	currentVer, err := svc.db.GetMaxVersion(ctx, repoID)
	if err != nil {
		return nil, err
	}
	if currentVer < 1 {
		return nil, apperrors.NewAppError(ErrorMissingSignedRole, "no active version found")
	}
	return svc.db.FindVersion(ctx, repoID, currentVer)
}

// createDefault creates the first data.SignedRootRole for data.Repo
func (svc *SignedContentService) createDefault(keys []data.RepoKey) (*data.RootRole, error) {
	if len(keys) == 0 {
		return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, "no keys found")
	}
	keyMap := make(map[data.KeyID]*encryption.SerializedKey)
	keyRolesMap := make(map[data.RoleType][]data.KeyID)
	for _, k := range keys {
		pkey, err := k.ToPublicKey()
		if err != nil {
			return nil, err
		}
		keyMap[k.KeyID] = pkey
		ids := make([]data.KeyID, 0)
		if s, ok := keyRolesMap[k.Role]; ok {
			ids = s
		}
		keyRolesMap[k.Role] = append(ids, k.KeyID)
	}
	return data.NewRootRole(keyMap, keyRolesMap, time.Now().Add(defaultRoleExpire)), nil
}

func (svc *SignedContentService) persistSignedPayload(ctx context.Context, role *data.RootRole) (*data.SignedRootRole, error) {
	// TODO: implement
	return nil, fmt.Errorf("not Implemented")
}

func (svc *SignedContentService) createAndPersist(ctx context.Context, repoID data.RepoID) (*data.SignedPayload[data.RootRole], error) {
	keys, err := svc.keySvc.GetRepoKeys(ctx, repoID)
	if err != nil {
		return nil, err
	}
	r, err := svc.createDefault(keys)
	if err != nil {
		return nil, err
	}
	sig, err := svc.persistSignedPayload(ctx, r)
	if err != nil {
		return nil, err
	}
	return &sig.Content, nil
}
