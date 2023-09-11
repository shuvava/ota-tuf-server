package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	intData "github.com/shuvava/ota-tuf-server/internal/data"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"github.com/shuvava/ota-tuf-server/pkg/encryption"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const keyRepoTableName = "tuf_keys"

type keyDTO struct {
	// Type is key type
	Type string `json:"keytype"`
	// Value is key value
	Value encryption.RawKey `json:"keyval"`
}

type repoKeyDTO struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	RepoID  string             `bson:"repoId"`
	Role    string             `bson:"role"`
	KeyID   string             `bson:"keyId"`
	Key     keyDTO             `bson:"key"`
	Created time.Time          `bson:"created"`
}

// RepoKeyMongoRepository implementations of db.KeyRepository for MongoDb repo
type RepoKeyMongoRepository struct {
	db   *intMongo.Db
	coll *mongo.Collection
	log  logger.Logger
	db.KeyRepository
}

// NewKeyMongoRepository creates new instance of RepoKeyMongoRepository
func NewKeyMongoRepository(logger logger.Logger, db *intMongo.Db) *RepoKeyMongoRepository {
	log := logger.SetOperation("TUFKeyRepo")
	return &RepoKeyMongoRepository{
		db:   db,
		coll: db.GetCollection(keyRepoTableName),
		log:  log,
	}
}

// Create persist new data.RepoKey in database
func (store *RepoKeyMongoRepository) Create(ctx context.Context, obj *data.RepoKey) error {
	log := store.log.
		SetOperation("Create").WithContext(ctx).
		WithField(intData.LogFieldRepoID, obj.RepoID).
		WithField(intData.LogFieldKeyID, obj.KeyID).
		WithField(intData.LogFieldRole, obj.Role)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Creating new SerializedKey")

	exists, err := store.Exists(ctx, obj.RepoID, obj.KeyID)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document(SerializedKey) with id='%s' role='%s' already exist in database", obj.KeyID, obj.Role)
		return apperrors.CreateErrorAndLogIt(log,
			ErrorRepoKeyErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	dto := toRepoKeyDTO(obj)
	_, err = store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.Info("SerializedKey created successful")
	} else {
		log.Warn("SerializedKey creation failed")
	}
	return err
}

// FindByKeyID returns data.RepoKey by keyID
func (store *RepoKeyMongoRepository) FindByKeyID(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (*data.RepoKey, error) {
	log := store.log.SetOperation("FindByKeyID").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID).
		WithField(intData.LogFieldKeyID, keyID)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Looking up RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	var dto repoKeyDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.Warn("RepoKey not found")
		}
		return nil, err
	}
	log.Debug("Object Found")
	model, err := toRepoKeyModel(dto)
	return model, err
}

// FindByRepoID returns data.RepoKey by repoId
func (store *RepoKeyMongoRepository) FindByRepoID(ctx context.Context, repoID data.RepoID) ([]*data.RepoKey, error) {
	log := store.log.SetOperation("FindByRepoID").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Looking up RepoKeys")

	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repoId", Value: repoID.String()}},
		},
	}}
	var docs []repoKeyDTO
	err := store.db.Find(ctx, store.coll, filter, &docs)
	if err != nil {
		log.Debug("Not Found")
		return nil, err
	}

	var res []*data.RepoKey
	for _, doc := range docs {
		obj, err := toRepoKeyModel(doc)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}

	log.WithField("Count", len(res)).
		Debug("Lookup completed successful")

	return res, nil
}

// Exists checks if data.RepoKey exists in database
func (store *RepoKeyMongoRepository) Exists(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (bool, error) {
	log := store.log.SetOperation("Exists").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID).
		WithField(intData.LogFieldKeyID, keyID)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Looking up RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// ExistsKeyRole checks if data.RepoKey with role exists in database
func (store *RepoKeyMongoRepository) ExistsKeyRole(ctx context.Context, repoID data.RepoID, role data.RoleType) (bool, error) {
	log := store.log.SetOperation("ExistsKeyRole").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID).
		WithField(intData.LogFieldRole, role)
	defer log.TrackFuncTime(time.Now())
	log.
		Debug("Looking up RepoKey")
	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repoId", Value: repoID}},
			bson.D{primitive.E{Key: "role", Value: role}},
		},
	}}
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// DeletePrivateKey removes private key
func (store *RepoKeyMongoRepository) DeletePrivateKey(ctx context.Context, repoID data.RepoID, keyID data.KeyID) error {
	log := store.log.SetOperation("DeletePrivateKey").
		WithContext(ctx).
		WithField(intData.LogFieldRepoID, repoID).
		WithField(intData.LogFieldKeyID, keyID)
	defer log.TrackFuncTime(time.Now())
	log.Debug("Deleting private key for RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	var dto repoKeyDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.
				Warn("RepoKey not found")
		}
		return err
	}
	dto.Key.Value.Private = nil
	err = store.db.ReplaceOne(ctx, store.coll, filter, &dto)
	if err != nil {
		log.WithField(intData.LogFieldError, err.Error()).
			Error("error on updating RepoKey")
	} else {
		log.Info("Private key deleted from RepoKey")
	}

	return err
}

// toRepoKeyDTO converts data.RepoKey to DTO
func toRepoKeyDTO(obj *data.RepoKey) repoKeyDTO {
	return repoKeyDTO{
		ID:      primitive.NewObjectID(),
		RepoID:  obj.RepoID.String(),
		Role:    string(obj.Role),
		KeyID:   obj.KeyID.String(),
		Created: obj.Created,
		Key: keyDTO{
			Type:  string(obj.Key.Type),
			Value: obj.Key.Value,
		},
	}
}

func toRepoKeyModel(dto repoKeyDTO) (*data.RepoKey, error) {
	repoID, err := data.RepoIDFromString(dto.RepoID)
	if err != nil {
		return nil, err
	}
	keyID, err := data.KeyIDFromString(dto.KeyID)
	if err != nil {
		return nil, err
	}
	key := data.RepoKey{
		RepoID:  repoID,
		Role:    data.RoleType(dto.Role),
		KeyID:   keyID,
		Created: dto.Created,
		Key: encryption.SerializedKey{
			Type:  encryption.KeyType(dto.Key.Type),
			Value: dto.Key.Value,
		},
	}
	return &key, nil
}

func getOneRepoKeyFilter(repoID data.RepoID, keyID data.KeyID) bson.D {
	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repoId", Value: repoID.String()}},
			bson.D{primitive.E{Key: "keyId", Value: keyID.String()}},
		},
	}}
	return filter
}
