package app

import (
	"context"
	"strings"

	intDb "github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	"github.com/shuvava/go-ota-svc-common/db"
	intCmnDb "github.com/shuvava/go-ota-svc-common/db/mongo"
)

func (s *Server) initDbService(ctx context.Context) {
	log := s.log.SetOperation("server-init-db")
	if s.svc.Db != nil {
		if err := s.svc.Db.Disconnect(ctx); err != nil {
			log.WithError(err).
				Fatal("Error on Db service distracting")
		}
	}
	switch db.Type(strings.ToLower(s.config.Db.Type)) {
	case db.MongoDb:
		mongoDB, err := intCmnDb.NewMongoDB(ctx, s.log, s.config.Db.ConnectionString)
		if err != nil {
			log.WithError(err).
				Fatal("Error on Db service creating")
		}
		// TODO Add db initialization (indexes creation) / schema migration functionality
		s.svc.Db = mongoDB
		s.svc.KeyRepo = intDb.NewKeyMongoRepository(s.log, mongoDB)
		s.svc.Repo = intDb.NewTUFRepoMongoRepository(s.log, mongoDB)
		s.svc.SignedContent = intDb.NewSignedContentMongoRepository(s.log, mongoDB)
	default:
		log.WithField("type", s.config.Db.Type).
			Fatal("Unsupported DB type")
	}
}

// create all application services
func (s *Server) initServices(ctx context.Context) {
	s.initDbService(ctx)
	s.svc.KeyRepoSvc = services.NewKeyRepositoryService(s.log, s.svc.KeyRepo)
	s.svc.SignedContentSvc = services.NewSignedContentService(s.log, s.svc.SignedContent)
	s.svc.RepoSvc = services.NewRepositoryService(s.log, s.svc.KeyRepoSvc, s.svc.SignedContentSvc, s.svc.Repo)
}
