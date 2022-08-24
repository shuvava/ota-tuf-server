package mongo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/data"

	cmndata "github.com/shuvava/go-ota-svc-common/data"
	intCmnDb "github.com/shuvava/go-ota-svc-common/db/mongo"
)

func TestMongoDB(t *testing.T) {
	var connStr = "mongodb://mongoadmin:secret@localhost:27017/test?authSource=admin"
	ctx := context.Background()
	log := logger.NewNopLogger()
	mongoDb, err := intCmnDb.NewMongoDB(ctx, log, connStr)
	if err != nil {
		t.Errorf("got %s, expected nil", err)
	}
	svc := mongo.NewTUFRepoMongoRepository(log, mongoDb)
	t.Run("should be created", func(t *testing.T) {
		ns := cmndata.NewNamespace(cmndata.NewCorrelationID().String())
		repo := data.Repo{
			Namespace: ns,
			RepoID:    data.NewRepoID(),
		}
		err = svc.Create(ctx, repo)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		repo.RepoID = data.NewRepoID()
		err = svc.Create(ctx, repo)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != mongo.ErrorRepoErrorDbAlreadyExist {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbConnection)
		}
	})
	t.Run("should return paginated result", func(t *testing.T) {
		var result []*data.Repo
		res, cnt, err := svc.List(ctx, 0, 2, nil)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if cnt == 0 {
			t.Errorf("got %d, expected > 0", cnt)
		}
		result = append(result, res...)
		_, _, err = svc.List(ctx, int64(len(result)), 2, nil)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
	})
}
