package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mdb "go.mongodb.org/mongo-driver/mongo"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	intMongo "github.com/shuvava/go-ota-svc-common/db/mongo"

	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/data"
)

const signedRepoTableName = "tuf_signed"

type repoSignedContentDTO struct {
	ID        primitive.ObjectID                `bson:"_id,omitempty"`
	RepoID    string                            `bson:"repo_id"`
	ExpiresAt time.Time                         `bson:"expires_at"`
	Version   int                               `bson:"version"`
	Content   data.SignedPayload[data.RootRole] `bson:"signed_payload"`
}

// SignedContentMongoRepository implementation of db.TufSignedContent for MongoDb repo
type SignedContentMongoRepository struct {
	db   *intMongo.Db
	coll *mdb.Collection
	log  logger.Logger
	db.TufSignedContent
}

// NewSignedContentMongoRepository creates new instance of SignedContentMongoRepository
func NewSignedContentMongoRepository(logger logger.Logger, db *intMongo.Db) *SignedContentMongoRepository {
	log := logger.SetOperation("SignedContentRepo")
	return &SignedContentMongoRepository{
		db:   db,
		coll: db.GetCollection(signedRepoTableName),
		log:  log,
	}
}

// Create persist new data.SignedRootRole in database
func (store *SignedContentMongoRepository) Create(ctx context.Context, obj *data.SignedRootRole) error {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", obj.RepoID).
		WithField("Version", obj.Version).
		Debug("Creating new SignedContent")

	exists, err := store.Exists(ctx, obj.RepoID, obj.Version)
	if err != nil {
		return err
	}
	if exists {
		err = fmt.Errorf("document(SerializedKey) with repo_id='%s' version='%d' already exist in database", obj.RepoID, obj.Version)
		return apperrors.CreateErrorAndLogIt(log,
			ErrorSignedContentErrorDbAlreadyExist,
			"Failed to add new DB record", err)
	}
	dto := toSignedContentDTO(obj)
	id, err := store.db.InsertOne(ctx, store.coll, dto)
	if err == nil {
		log.WithField("RepoID", obj.RepoID).
			WithField("Version", obj.Version).
			WithField("ObjectID", id).
			Info("SignedContent created successful")
	} else {
		log.WithField("RepoID", obj.RepoID).
			WithField("Version", obj.Version).
			Warn("SignedContent creation failed")
	}
	return err
}

// Exists checks if data.SignedRootRole exists in database
func (store *SignedContentMongoRepository) Exists(ctx context.Context, repoID data.RepoID, ver int) (bool, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		WithField("Version", ver).
		Debug("Checkin existence of the SignedContent")
	filter := getOneSignedContentFilter(repoID, ver)
	cnt, err := store.db.Count(ctx, store.coll, filter)
	if err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// FindVersion returns data.SignedRootRole with Version equal to ver parameter
func (store *SignedContentMongoRepository) FindVersion(ctx context.Context, repoID data.RepoID, ver int) (*data.SignedRootRole, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		WithField("Version", ver).
		Debug("Looking up SignedContent")
	filter := getOneSignedContentFilter(repoID, ver)
	var dto repoSignedContentDTO
	err := store.db.GetOne(ctx, store.coll, filter, &dto)
	if err != nil {
		var typedErr apperrors.AppError
		if errors.As(err, &typedErr) && typedErr.ErrorCode == apperrors.ErrorDbNoDocumentFound {
			log.WithField("RepoID", repoID).
				WithField("Version", ver).
				Warn("SignedContent not found")
		}
		return nil, err
	}
	log.WithField("RepoID", repoID).
		WithField("Version", ver).
		Debug("Object Found")
	return toSignedContentModel(dto)
}

// GetMaxVersion returns current repo version
func (store *SignedContentMongoRepository) GetMaxVersion(ctx context.Context, repoID data.RepoID) (int, error) {
	log := store.log.WithContext(ctx)
	defer log.TrackFuncTime(time.Now())
	log.WithField("RepoID", repoID).
		Debug("Getting max version of SignedContent")
	pipeline := make([]bson.M, 0)
	groupStage := bson.M{
		"$group": bson.M{
			"_id":        nil,
			"maxVersion": bson.M{"$max": "$version"},
		},
	}
	matchStage := bson.M{
		"$match": bson.M{
			"repo_id": repoID.String(),
		},
	}
	pipeline = append(pipeline, matchStage, groupStage)
	var res []intMongo.DBResult
	if err := store.db.Aggregate(ctx, store.coll, pipeline, nil, &res); err != nil {
		return -1, err
	}
	if len(res) < 1 {
		return 0, apperrors.NewAppError(apperrors.ErrorDbOperation, "unexpected aggregation result")
	}

	return int(res[0]["maxVersion"].(int32)), nil
}

func getOneSignedContentFilter(repoID data.RepoID, ver int) bson.D {
	return bson.D{primitive.E{
		Key: "$and",
		Value: bson.A{
			bson.D{primitive.E{Key: "repo_id", Value: repoID.String()}},
			bson.D{primitive.E{Key: "version", Value: ver}},
		},
	}}
}

// toSignedContentDTO converts data.SignedRootRole to repoSignedContentDTO
func toSignedContentDTO(obj *data.SignedRootRole) repoSignedContentDTO {
	return repoSignedContentDTO{
		ID:        primitive.NewObjectID(),
		RepoID:    obj.RepoID.String(),
		ExpiresAt: obj.ExpiresAt,
		Version:   obj.Version,
		Content:   obj.Content,
	}
}

func toSignedContentModel(dto repoSignedContentDTO) (*data.SignedRootRole, error) {
	repoID, err := data.RepoIDFromString(dto.RepoID)
	if err != nil {
		return nil, err
	}
	return &data.SignedRootRole{
		RepoID:    repoID,
		ExpiresAt: dto.ExpiresAt,
		Version:   dto.Version,
		Content:   dto.Content,
	}, nil
}
