package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

const objectTableName = "tuf_keys"

type keyDTO struct {
	// Type is key type
	Type string `json:"keytype"`
	// Value is key value
	Value json.RawMessage `json:"keyval"`
}

type repoKeyDTO struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	RepoID string             `json:"repo_id"`
	Role   string             `json:"role"`
	KeyID  string             `json:"key_id"`
	Key    keyDTO             `json:"key"`
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
	log := logger.SetOperation("KeyRepo")
	return &RepoKeyMongoRepository{
		db:   db,
		coll: db.GetCollection(objectTableName),
		log:  log,
	}
}

// Create persist new data.Object in database
func (store *RepoKeyMongoRepository) Create(ctx context.Context, obj data.RepoKey) error {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", obj.RepoID).
		WithField("KeyID", obj.KeyID).
		WithField("Role", obj.Role).
		Debug("Creating new Key")

	dto := toDTO(obj)
	exists, err := store.Exists(ctx, obj.RepoID, obj.KeyID)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document(Key) with id='%s' role='%s' already exist in database", obj.KeyID, obj.Role)
		return apperrors.CreateErrorAndLogIt(log,
			ErrorRepoKeyErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	_, err = store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.WithField("RepoID", obj.RepoID).
			WithField("KeyID", obj.KeyID).
			WithField("Role", obj.Role).
			Info("Key created successful")
	} else {
		log.WithField("RepoID", obj.RepoID).
			WithField("KeyID", obj.KeyID).
			WithField("Role", obj.Role).
			Warn("Key creation failed")
	}
	return err
}

// FindByKeyID returns data.RepoKey by keyID
func (store *RepoKeyMongoRepository) FindByKeyID(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (*data.RepoKey, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		WithField("KeyID", keyID).
		Debug("Looking up RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	var dto repoKeyDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("RepoID", repoID).
				WithField("KeyID", keyID).
				Warn("RepoKey not found")
		}
		return nil, err
	}
	log.WithField("RepoID", repoID).
		WithField("KeyID", keyID).
		Debug("Object Found")
	model, err := toModel(dto)
	return &model, err
}

// FindByRepoId returns data.RepoKey by repoId
func (store *RepoKeyMongoRepository) FindByRepoId(ctx context.Context, repoID data.RepoID) ([]data.RepoKey, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		Debug("Looking up RepoKeys")

	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repo_id", Value: repoID.String()}},
		},
	}}
	var docs []repoKeyDTO
	err := store.db.Find(ctx, store.coll, filter, &docs)
	if err != nil {
		log.WithField("RepoID", repoID).
			Debug("Not Found")
		return nil, err
	}

	var res []data.RepoKey
	for _, doc := range docs {
		obj, err := toModel(doc)
		if err != nil {
			return nil, err
		}
		res = append(res, obj)
	}

	log.WithField("RepoID", repoID).
		WithField("Count", len(res)).
		Debug("Lookup completed successful")

	return res, nil
}

// Exists checks if data.RepoKey exists in database
func (store *RepoKeyMongoRepository) Exists(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (bool, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		WithField("KeyID", keyID).
		Debug("Looking up RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// toDTO converts data.RepoKey to DTO
func toDTO(obj data.RepoKey) repoKeyDTO {
	return repoKeyDTO{
		ID:     primitive.NewObjectID(),
		RepoID: obj.RepoID.String(),
		Role:   string(obj.Role),
		KeyID:  obj.KeyID.String(),
		Key: keyDTO{
			Type:  string(obj.Key.Type),
			Value: obj.Key.Value,
		},
	}
}

func toModel(dto repoKeyDTO) (data.RepoKey, error) {
	repoID, err := data.RepoIDFromString(dto.RepoID)
	if err != nil {
		return data.RepoKey{}, err
	}
	keyID, err := data.KeyIDFromString(dto.KeyID)
	if err != nil {
		return data.RepoKey{}, err
	}
	return data.RepoKey{
		RepoID: repoID,
		Role:   data.RoleType(dto.Role),
		KeyID:  keyID,
		Key: data.Key{
			Type:  data.KeyType(dto.Key.Type),
			Value: dto.Key.Value,
		},
	}, nil
}

func getOneRepoKeyFilter(repoID data.RepoID, keyID data.KeyID) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repo_id", Value: repoID}},
			bson.D{primitive.E{Key: "key_id", Value: keyID}},
		},
	}}
}
