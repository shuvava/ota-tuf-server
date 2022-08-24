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
	defaultRoleExpire  = time.Hour * 24 * 365
	firstVersionNumber = 0
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
			return svc.createAndPersist(ctx, repoID, firstVersionNumber)
		} else {
			return nil, err
		}
	}
	if sig.ExpiresAt.Add(time.Hour).After(time.Now().UTC()) {
		return svc.createAndPersist(ctx, repoID, sig.Version)
	}

	// TODO: implement
	return sig.Content, nil
}

// getCurrent returns the latest created data.SignedRootRole for data.Repo
func (svc *SignedContentService) getCurrent(ctx context.Context, repoID data.RepoID) (*data.SignedRootRole, error) {
	currentVer, err := svc.db.GetMaxVersion(ctx, repoID)
	if err != nil {
		return nil, err
	}
	if currentVer <= firstVersionNumber {
		return nil, apperrors.NewAppError(ErrorMissingSignedRole, "no active version found")
	}
	return svc.db.FindVersion(ctx, repoID, currentVer)
}

// createNewVersion creates the next version of data.RootRole for data.Repo
func (svc *SignedContentService) createNewVersion(pervVer uint, keys []data.RepoKey) (*data.RootRole, error) {
	keyMap, keyRolesMap, err := createKeyMaps(keys)
	if err != nil {
		return nil, err
	}
	return data.NewRootRole(keyMap, keyRolesMap, pervVer+1, time.Now().Add(defaultRoleExpire)), nil
}

// signRootRole sings data.RootRole for data.RepoID
func (svc *SignedContentService) signRootRole(repoID data.RepoID, role *data.RootRole, keys map[data.KeyID]data.RepoKey) (*data.SignedRootRole, error) {
	rootKeyRole, ok := role.KeyRoles[data.RoleTypeRoot]
	if !ok || len(rootKeyRole.Keys) == 0 {
		return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, "failed to find key with root role")
	}
	rootKeys := make([]data.RepoKey, len(rootKeyRole.Keys))
	for i, keyID := range rootKeyRole.Keys {
		if key, ok := keys[keyID]; ok {
			rootKeys[i] = key
		} else {
			return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, fmt.Sprintf("failed to find key with ID %s", keyID))
		}
	}
	sig, err := role.Sign(rootKeys)
	if err != nil {
		return nil, err
	}
	signedRole := &data.SignedRootRole{
		RepoID:    repoID,
		ExpiresAt: role.ExpiresAt,
		Version:   role.Version,
		Content:   sig,
	}

	return signedRole, nil
}

// createAndPersist creates data.SignedRootRole with version 1 and persist in data storage
func (svc *SignedContentService) createAndPersist(ctx context.Context, repoID data.RepoID, pervVer uint) (*data.SignedPayload[data.RootRole], error) {
	keys, err := svc.keySvc.GetRepoKeys(ctx, repoID)
	if err != nil {
		return nil, err
	}
	r, err := svc.createNewVersion(pervVer, keys)
	if err != nil {
		return nil, err
	}
	keyMap := createMap(keys)
	sig, err := svc.signRootRole(repoID, r, keyMap)
	if err != nil {
		return nil, err
	}
	if err = svc.db.Create(ctx, sig); err != nil {
		return nil, err
	}
	return sig.Content, nil
}

func createMap(keys []data.RepoKey) map[data.KeyID]data.RepoKey {
	m := map[data.KeyID]data.RepoKey{}
	for _, key := range keys {
		m[key.KeyID] = key
	}
	return m
}

func createKeyMaps(keys []data.RepoKey) (map[data.KeyID]*encryption.SerializedKey, map[data.RoleType][]data.KeyID, error) {
	if len(keys) == 0 {
		return nil, nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, "no keys found")
	}
	keyMap := make(map[data.KeyID]*encryption.SerializedKey)
	keyRolesMap := make(map[data.RoleType][]data.KeyID)
	for _, k := range keys {
		pkey, err := k.ToPublicKey()
		if err != nil {
			return nil, nil, err
		}
		keyMap[k.KeyID] = pkey
		ids := make([]data.KeyID, 0)
		if s, ok := keyRolesMap[k.Role]; ok {
			ids = s
		}
		keyRolesMap[k.Role] = append(ids, k.KeyID)
	}
	return keyMap, keyRolesMap, nil
}
