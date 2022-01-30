package app

import (
	"context"
	"strings"

	intDb "github.com/shuvava/ota-tuf-server/internal/db/mongo"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	"github.com/shuvava/go-ota-svc-common/db"
	intCmnDb "github.com/shuvava/go-ota-svc-common/db/mongo"
)

func (s *Server) initDbService() {
	log := s.log.SetOperation("server-init-db")
	if s.svc.Db != nil {
		if err := s.svc.Db.Disconnect(context.Background()); err != nil {
			log.WithError(err).
				Fatal("Error on Db service distracting")
		}
	}
	switch db.Type(strings.ToLower(s.config.Db.Type)) {
	case db.MongoDb:
		mongoDB, err := intCmnDb.NewMongoDB(context.Background(), s.log, s.config.Db.ConnectionString)
		if err != nil {
			log.WithError(err).
				Fatal("Error on Db service creating")
		}
		s.svc.Db = mongoDB
		s.svc.KeyRepo = intDb.NewKeyMongoRepository(s.log, mongoDB)
	default:
		log.WithField("type", s.config.Db.Type).
			Fatal("Unsupported mongoDB type")
	}
}

// create all application services
func (s *Server) initServices() {
	s.initDbService()
	s.svc.KeySvc = services.NewRepositoryService(s.log, s.svc.KeyRepo)
}
