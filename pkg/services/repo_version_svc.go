package services

import (
	"context"
	"fmt"
	"time"

	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"

	"github.com/shuvava/go-logging/logger"
)

const (
	defaultRoleExpire  = time.Hour * 24 * 365
	firstVersionNumber = 0
)

// RepoVersionService controls TUF SerializedKey roles life cycle
type RepoVersionService struct {
	log logger.Logger
	db  db.TufSignedContent
}

// NewRepoVersionService creates new instance of services.RepoVersionService
func NewRepoVersionService(l logger.Logger, db db.TufSignedContent) *RepoVersionService {
	log := l.SetArea("signed-content-service")
	return &RepoVersionService{
		log: log,
		db:  db,
	}
}

// GetCurrentVersion returns current signature for data.RepoID
func (svc *RepoVersionService) GetCurrentVersion(ctx context.Context, repoID data.RepoID) (*data.SignedRootRole, error) {
	currentVer, err := svc.db.GetMaxVersion(ctx, repoID)
	if err != nil {
		return nil, err
	}
	if currentVer <= firstVersionNumber {
		return nil, apperrors.NewAppError(ErrorMissingSignedRole, "no active version found")
	}
	return svc.db.FindVersion(ctx, repoID, currentVer)
}

// GetVersion returns signature of data.RepoID for version
func (svc *RepoVersionService) GetVersion(ctx context.Context, repoID data.RepoID, version uint) (*data.SignedRootRole, error) {
	log := svc.log.SetOperation("GetVersion").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	exist, err := svc.db.Exists(ctx, repoID, version)
	if err != nil {
		return nil, err
	}
	if !exist {
		return nil, apperrors.NewAppError(ErrorMissingSignedRole, "version not found")
	}

	return svc.db.FindVersion(ctx, repoID, version)
}

// CreateFirstVersion creates the first data.SignedRootRole for the data.RepoID
func (svc *RepoVersionService) CreateFirstVersion(ctx context.Context, repoID data.RepoID, keys []*data.RepoKey, threshold uint) error {
	log := svc.log.SetOperation("CreateFirstVersion").WithContext(ctx)
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

// UpdateRepoVersion update signature of the data.Repo
func (svc *RepoVersionService) UpdateRepoVersion(ctx context.Context, sig *data.SignedRootRole, keys []*data.RepoKey) (*data.SignedRootRole, error) {
	log := svc.log.SetOperation("UpdateRepoVersion").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	return svc.createAndPersist(ctx, sig.RepoID, keys, sig.Version, sig.Threshold)
}

// createAndPersist creates data.SignedRootRole with version 1 and persist in data storage
func (svc *RepoVersionService) createAndPersist(ctx context.Context, repoID data.RepoID, keys []*data.RepoKey, pervVer uint, threshold uint) (*data.SignedRootRole, error) {
	r, err := createNewVersion(pervVer, keys, threshold)
	if err != nil {
		return nil, err
	}
	keyMap := createMap(keys)
	sig, err := signRepo(repoID, r, keyMap, threshold)
	if err != nil {
		return nil, err
	}
	if err = svc.db.Create(ctx, sig); err != nil {
		return nil, err
	}
	return sig, nil
}

// createNewVersion creates the next version of data.RepoSigned for data.Repo
func createNewVersion(pervVer uint, keys []*data.RepoKey, threshold uint) (*data.RepoSigned, error) {
	keyMap, keyRolesMap, err := createKeyMaps(keys)
	if err != nil {
		return nil, err
	}
	return data.NewRepoSigned(keyMap, keyRolesMap, pervVer+1, time.Now().Add(defaultRoleExpire), threshold), nil
}

// signRepo sings data.RepoSigned for data.RepoID
func signRepo(repoID data.RepoID, repoSigned *data.RepoSigned, keys map[data.KeyID]*data.RepoKey, threshold uint) (*data.SignedRootRole, error) {
	rootKeyRole, ok := repoSigned.KeyRoles[data.RoleTypeRoot]
	if !ok || len(rootKeyRole.Keys) == 0 {
		return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, "failed to find key with root repoSigned")
	}
	rootKeys := make([]*data.RepoKey, len(rootKeyRole.Keys))
	for i, keyID := range rootKeyRole.Keys {
		if key, ok := keys[keyID]; ok {
			rootKeys[i] = key
		} else {
			return nil, apperrors.NewAppError(ErrorSvcSignedContentKeyNotFound, fmt.Sprintf("failed to find key with ID %s", keyID))
		}
	}
	sig, err := repoSigned.Sign(rootKeys)
	if err != nil {
		return nil, err
	}
	signedRole := &data.SignedRootRole{
		RepoID:    repoID,
		ExpiresAt: repoSigned.ExpiresAt,
		Version:   repoSigned.Version,
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
