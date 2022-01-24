package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/shuvava/go-logging/logger"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

// RepoKeyMongoRepository implementations of db.RepoKeyRepository for MongoDb repo
type RepoKeyMongoRepository struct {
	db   *intMongo.Db
	coll *mongo.Collection
	log  logger.Logger
	db.RepoKeyRepository
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
	log.WithField("RepoId", obj.RepoID).
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
		log.WithField("RepoId", obj.RepoID).
			WithField("KeyID", obj.KeyID).
			WithField("Role", obj.Role).
			Debug("Key created successful")
	} else {
		log.WithField("RepoId", obj.RepoID).
			WithField("KeyID", obj.KeyID).
			WithField("Role", obj.Role).
			Warn("Key creation failed")
	}
	return err
}

// FindByKeyID returns data.RepoKey by keyID
func (store *RepoKeyMongoRepository) FindByKeyID(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (*data.RepoKey, error) {
	log := store.log.WithContext(ctx)
	log.WithField("RepoId", repoID).
		WithField("KeyID", keyID).
		Debug("Looking up RepoKey")
	filter := getOneRepoKeyFilter(repoID, keyID)
	var dto repoKeyDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("RepoId", repoID).
				WithField("KeyID", keyID).
				Warn("RepoKey not found")
		}
		return nil, err
	}
	log.WithField("RepoId", repoID).
		WithField("KeyID", keyID).
		Debug("Object Found")
	model := toModel(dto)
	return &model, nil
}

// FindByRepoId returns data.RepoKey by repoId
func (store *RepoKeyMongoRepository) FindByRepoId(ctx context.Context, repoID data.RepoID) ([]data.RepoKey, error) {
	log := store.log.WithContext(ctx)
	log.WithField("RepoId", repoID).
		Debug("Looking up RepoKeys")

	filter := bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repo_id", Value: string(repoID)}},
		},
	}}
	var docs []repoKeyDTO
	err := store.db.Find(ctx, store.coll, filter, &docs)
	if err != nil {
		log.WithField("RepoId", repoID).
			Debug("Not Found")
		return nil, err
	}

	var res []data.RepoKey
	for _, doc := range docs {
		obj := toModel(doc)
		res = append(res, obj)
	}

	log.WithField("RepoId", repoID).
		WithField("Count", len(res)).
		Debug("Lookup completed successful")

	return res, nil
}

// Exists checks if data.RepoKey exists in database
func (store *RepoKeyMongoRepository) Exists(ctx context.Context, repoID data.RepoID, keyID data.KeyID) (bool, error) {
	log := store.log.WithContext(ctx)
	log.WithField("RepoId", repoID).
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
		RepoID: string(obj.RepoID),
		Role:   string(obj.Role),
		KeyID:  string(obj.KeyID),
		Key: keyDTO{
			Type:  string(obj.Key.Type),
			Value: obj.Key.Value,
		},
	}
}

func toModel(dto repoKeyDTO) data.RepoKey {
	return data.RepoKey{
		RepoID: data.RepoID(dto.RepoID),
		Role:   data.RoleType(dto.Role),
		KeyID:  data.KeyID(dto.KeyID),
		Key: data.Key{
			Type:  data.KeyType(dto.Key.Type),
			Value: dto.Key.Value,
		},
	}
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
