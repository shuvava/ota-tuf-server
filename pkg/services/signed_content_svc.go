package services

import (
	"context"
	"fmt"
	"time"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"
)

const (
	defaultRoleExpire  = time.Hour //* 24 * 365
	firstVersionNumber = 0
)

// SignedContentService controls TUF SerializedKey roles life cycle
type SignedContentService struct {
	log logger.Logger
	db  db.TufSignedContent
}

// NewSignedContentService creates new instance of services.SignedContentService
func NewSignedContentService(l logger.Logger, db db.TufSignedContent) *SignedContentService {
	log := l.SetArea("signed-content-service")
	return &SignedContentService{
		log: log,
		db:  db,
	}
}

// GetCurrentSignature returns current signature for data.RepoID
func (svc *SignedContentService) GetCurrentSignature(ctx context.Context, repoID data.RepoID) (*data.SignedRootRole, error) {
	currentVer, err := svc.db.GetMaxVersion(ctx, repoID)
	if err != nil {
		return nil, err
	}
	if currentVer <= firstVersionNumber {
		return nil, apperrors.NewAppError(ErrorMissingSignedRole, "no active version found")
	}
	return svc.db.FindVersion(ctx, repoID, currentVer)
}

// CreateNewSignature creates the first data.SignedRootRole for the data.RepoID
func (svc *SignedContentService) CreateNewSignature(ctx context.Context, repoID data.RepoID, keys []*data.RepoKey, threshold uint) error {
	log := svc.log.SetOperation("CreateNewSignature").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	b, err := svc.db.Exists(ctx, repoID, firstVersionNumber)
	if err != nil {
		return err
	}
	if b {
		return apperrors.NewAppError(ErrorVersionAlreadyExist, fmt.Sprintf("repo with ID %s and version %d already exist", repoID, firstVersionNumber))
	}
	_, err = svc.createAndPersist(ctx, repoID, keys, firstVersionNumber, threshold)
	return err
}

// UpdateSignature update signature of the data.Repo
func (svc *SignedContentService) UpdateSignature(ctx context.Context, sig *data.SignedRootRole, keys []*data.RepoKey) (*data.SignedRootRole, error) {
	log := svc.log.SetOperation("UpdateSignature").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	//currentVer, err := svc.db.GetMaxVersion(ctx, repoID)
	//if err != nil {
	//	return nil, err
	//}
	//if currentVer <= firstVersionNumber {
	//	return nil, apperrors.NewAppError(ErrorMissingSignedRole, "no active version found")
	//}
	//sig, err := svc.db.FindVersion(ctx, repoID, currentVer)
	//if err != nil {
	//	return nil, err
	//}
	//if sig.ExpiresAt.Before(time.Now().Add(time.Hour).UTC()) {
	//	// should sign new version
	//}
	return svc.createAndPersist(ctx, sig.RepoID, keys, sig.Version, sig.Threshold)
}

// createAndPersist creates data.SignedRootRole with version 1 and persist in data storage
func (svc *SignedContentService) createAndPersist(ctx context.Context, repoID data.RepoID, keys []*data.RepoKey, pervVer uint, threshold uint) (*data.SignedRootRole, error) {
	r, err := createNewVersion(pervVer, keys, threshold)
	if err != nil {
		return nil, err
	}
	keyMap := createMap(keys)
	sig, err := signRootRole(repoID, r, keyMap, threshold)
	if err != nil {
		return nil, err
	}
	if err = svc.db.Create(ctx, sig); err != nil {
		return nil, err
	}
	return sig, nil
}

// createNewVersion creates the next version of data.RootRole for data.Repo
func createNewVersion(pervVer uint, keys []*data.RepoKey, threshold uint) (*data.RootRole, error) {
	keyMap, keyRolesMap, err := createKeyMaps(keys)
	if err != nil {
		return nil, err
	}
	return data.NewRootRole(keyMap, keyRolesMap, pervVer+1, time.Now().Add(defaultRoleExpire), threshold), nil
}

// signRootRole sings data.RootRole for data.RepoID
func signRootRole(repoID data.RepoID, role *data.RootRole, keys map[data.KeyID]*data.RepoKey, threshold uint) (*data.SignedRootRole, error) {
	rootKeyRole, ok := role.KeyRoles[data.RoleTypeRoot]
	if !ok || len(rootKeyRole.Keys) == 0 {
		return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, "failed to find key with root role")
	}
	rootKeys := make([]*data.RepoKey, len(rootKeyRole.Keys))
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
		Threshold: threshold,
	}

	return signedRole, nil
}

func createMap(keys []*data.RepoKey) map[data.KeyID]*data.RepoKey {
	m := map[data.KeyID]*data.RepoKey{}
	for _, key := range keys {
		m[key.KeyID] = key
	}
	return m
}

func createKeyMaps(keys []*data.RepoKey) (map[data.KeyID]*encryption.SerializedKey, map[data.RoleType][]data.KeyID, error) {
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
