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
	cmndata "github.com/shuvava/go-ota-svc-common/data"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const repoTableName = "tuf_repos"

type repoDTO struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Namespace string             `bson:"namespace"`
	RepoID    string             `bson:"repoId"`
	KeyType   string             `json:"keyType"`
}

// TUFRepoMongoRepository implementations of db.TufRepoRepository for MongoDb repo
type TUFRepoMongoRepository struct {
	db   *intMongo.Db
	coll *mongo.Collection
	log  logger.Logger

	db.TufRepoRepository
}

// NewTUFRepoMongoRepository creates new instance of TUFRepoMongoRepository
func NewTUFRepoMongoRepository(logger logger.Logger, db *intMongo.Db) *TUFRepoMongoRepository {
	log := logger.SetArea("TUFRepo")
	return &TUFRepoMongoRepository{
		db:   db,
		coll: db.GetCollection(repoTableName),
		log:  log,
	}
}

// Create persist new data.Repo in database
func (store *TUFRepoMongoRepository) Create(ctx context.Context, obj *data.Repo) error {
	log := store.log.
		SetOperation("Create").
		WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField(intData.LogFieldRepoID, obj.RepoID).
		WithField(intData.LogFieldNamespace, obj.Namespace).
		Debug("Creating new Repo")

	dto := toRepoDTO(obj)
	exists, err := store.Exists(ctx, obj.Namespace)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document(Repo) with id='%s' namespace='%s' already exist in database", obj.RepoID, obj.Namespace)
		return apperrors.CreateErrorAndLogIt(log,
			ErrorRepoErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	_, err = store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.WithField(intData.LogFieldRepoID, obj.RepoID).
			WithField(intData.LogFieldNamespace, obj.Namespace).
			Info("Repo created successful")
	} else {
		log.WithField(intData.LogFieldRepoID, obj.RepoID).
			WithField(intData.LogFieldNamespace, obj.Namespace).
			Warn("Repo creation failed")
	}
	return err
}

// FindByNamespace returns data.Repo by Namespace
func (store *TUFRepoMongoRepository) FindByNamespace(ctx context.Context, ns cmndata.Namespace) (*data.Repo, error) {
	log := store.log.SetOperation("FindByNamespace").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField(intData.LogFieldNamespace, ns).
		Debug("Looking up Repo by Namespace")
	filter := getOneRepoFilterByNamespace(ns)
	var dto repoDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField(intData.LogFieldNamespace, ns).
				Warn("Repo not found")
		}
		return nil, err
	}
	log.WithField(intData.LogFieldNamespace, ns).
		Debug("Repo Found")
	model, err := toRepoModel(dto)
	return &model, err
}

// FindByID returns data.Repo by RepoID
func (store *TUFRepoMongoRepository) FindByID(ctx context.Context, repoID data.RepoID) (*data.Repo, error) {
	log := store.log.SetOperation("FindByID").WithContext(ctx)
	log.WithField(intData.LogFieldRepoID, repoID).
		Debug("Looking up Repo by " + intData.LogFieldRepoID)
	filter := getOneRepoFilterByRepoID(repoID)
	var dto repoDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField(intData.LogFieldRepoID, repoID).
				Warn("Repo not found")
		}
		return nil, err
	}
	log.WithField(intData.LogFieldRepoID, repoID).
		Debug("Repo Found")
	model, err := toRepoModel(dto)
	return &model, err
}

// Exists checks if data.RepoKey exists in database
func (store *TUFRepoMongoRepository) Exists(ctx context.Context, ns cmndata.Namespace) (bool, error) {
	log := store.log.SetOperation("Exists").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField(intData.LogFieldNamespace, ns).
		Debug("Looking up Repo by Namespace")
	filter := getOneRepoFilterByNamespace(ns)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// List returns all data.Repo
func (store *TUFRepoMongoRepository) List(ctx context.Context, skip, limit int64, sortFiled *string) ([]*data.Repo, int64, error) {
	log := store.log.SetOperation("List").WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("Skip", skip).
		WithField("Limit", limit).
		Debug("Getting all Repos")
	var facetData []bson.M
	if sortFiled != nil && *sortFiled != "" {
		facetData = append(facetData, bson.M{"$sort": bson.E{Key: *sortFiled, Value: 1}})
	}
	facetData = append(facetData, bson.M{"$skip": skip})
	facetData = append(facetData, bson.M{"$limit": limit})
	facet := bson.M{"$facet": bson.M{
		"data":  facetData,
		"total": []bson.M{{"$count": "count"}},
	}}
	diskUse := true
	opt := &options.AggregateOptions{
		AllowDiskUse: &diskUse,
	}
	var dbRes []*paginatedResponse
	if err := store.db.Aggregate(ctx, store.coll, []bson.M{facet}, opt, &dbRes); err != nil {
		return nil, 0, err
	}
	if len(dbRes) < 1 {
		return nil, 0, apperrors.NewAppError(apperrors.ErrorDbOperation, "unexpected aggregation result")
	}
	var res []*data.Repo
	for _, dto := range dbRes[0].Data {
		if model, err := toRepoModel(dto); err == nil {
			res = append(res, &model)
		}
	}
	return res, dbRes[0].Total[0].Count, nil
}

func getOneRepoFilterByNamespace(ns cmndata.Namespace) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "namespace", Value: string(ns)}},
		},
	}}
}

func getOneRepoFilterByRepoID(repoID data.RepoID) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repoId", Value: repoID.String()}},
		},
	}}
}

// toRepoKeyDTO converts data.RepoKey to DTO
func toRepoDTO(obj *data.Repo) repoDTO {
	return repoDTO{
		ID:        primitive.NewObjectID(),
		RepoID:    obj.RepoID.String(),
		Namespace: string(obj.Namespace),
		KeyType:   string(obj.KeyType),
	}
}

func toRepoModel(dto repoDTO) (data.Repo, error) {
	repoID, err := data.RepoIDFromString(dto.RepoID)
	if err != nil {
		return data.Repo{}, err
	}
	return data.Repo{
		RepoID:    repoID,
		Namespace: cmndata.Namespace(dto.Namespace),
		KeyType:   encryption.KeyTypeFromString(dto.KeyType),
	}, nil
}

type paginatedResponse struct {
	Total []struct {
		Count int64 `json:"count"`
	} `json:"total"`
	Data []repoDTO `json:"data"`
}
