package mongo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/shuvava/go-logging/logger"
	"github.com/shuvava/go-ota-svc-common/apperrors"
	"github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/data"

	intCmnDb "github.com/shuvava/go-ota-svc-common/db/mongo"
)

func TestSignedContentMongoRepository(t *testing.T) {
	var connStr = "mongodb://mongoadmin:secret@localhost:27017/test?authSource=admin"
	ctx := context.Background()
	log := logger.NewNopLogger()
	mongoDb, err := intCmnDb.NewMongoDB(ctx, log, connStr)
	if err != nil {
		t.Errorf("got %s, expected nil", err)
	}
	svc := mongo.NewSignedContentMongoRepository(log, mongoDb)
	repoID := data.NewRepoID()
	content := "this is test string blablablablablablabla"
	sig := &data.ClientSignature{}
	sig.Value = content
	t.Run("should be created", func(t *testing.T) {
		obj := &data.SignedRootRole{
			RepoID:  repoID,
			Version: 1,
			Content: &data.SignedPayload[data.RootRole]{
				Signatures: []*data.ClientSignature{sig},
			},
			ExpiresAt: time.Now().Add(time.Hour * 3),
		}
		err = svc.Create(ctx, obj)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		err = svc.Create(ctx, obj)
		var typedErr apperrors.AppError
		if err == nil || errors.As(err, &typedErr) && typedErr.ErrorCode != mongo.ErrorSignedContentErrorDbAlreadyExist {
			t.Errorf("got %s, expected %s", err, apperrors.ErrorDbConnection)
		}
	})
	t.Run("should be read", func(t *testing.T) {
		obj := &data.SignedRootRole{
			RepoID:  repoID,
			Version: 2,
			Content: &data.SignedPayload[data.RootRole]{
				Signatures: []*data.ClientSignature{sig},
			},
			ExpiresAt: time.Now().Add(time.Hour * 3),
		}
		if err = svc.Create(ctx, obj); err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		obj2, err := svc.FindVersion(ctx, repoID, 2)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if len(obj2.Content.Signatures) != len(obj.Content.Signatures) ||
			len(obj2.Content.Signatures[0].Value) != len(obj.Content.Signatures[0].Value) {
			t.Errorf("content does not match")
		}
	})
	t.Run("should return correct max version", func(t *testing.T) {
		obj := &data.SignedRootRole{
			RepoID:    repoID,
			Version:   100,
			Content:   &data.SignedPayload[data.RootRole]{},
			ExpiresAt: time.Now().Add(time.Hour * 3),
		}
		if err = svc.Create(ctx, obj); err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		ver, err := svc.GetMaxVersion(ctx, repoID)
		if err != nil {
			t.Errorf("got %s, expected nil", err)
		}
		if ver != obj.Version {
			t.Errorf("Incorrect next version: %d should be %d", ver, obj.Version)
		}
	})
	t.Run("should return error if repo does not exist", func(t *testing.T) {
		r := data.NewRepoID()
		ver, err := svc.GetMaxVersion(ctx, r)
		if err == nil {
			t.Errorf("got %s, expected nil", err)
		}
		if ver != 0 {
			t.Errorf("Incorrect next version: %d should be %d", ver, 0)
		}
	})
}
