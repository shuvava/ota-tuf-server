package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/shuvava/ota-tuf-server/internal/config"
	"github.com/shuvava/ota-tuf-server/internal/db"
	"github.com/shuvava/ota-tuf-server/pkg/services"

	"github.com/shuvava/go-logging/logger"
	intCmnDb "github.com/shuvava/go-ota-svc-common/db"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// Server is main application servers
type Server struct {
	Echo   *echo.Echo
	log    logger.Logger
	config *config.AppConfig
	mu     sync.Mutex
	svc    struct {
		Db               intCmnDb.BaseRepository
		KeyRepo          db.KeyRepository
		Repo             db.TufRepoRepository
		SignedContent    db.TufSignedContent
		RepoSvc          *services.RepositoryService
		KeyRepoSvc       *services.KeyRepositoryService
		SignedContentSvc *services.RepoVersionService
	}
}

// NewServer creates new Server instance
func NewServer(logger logger.Logger) *Server {
	s := &Server{
		log: logger,
	}

	s.initWebServer()
	s.initConfig()

	return s
}

// initConfig load app config
func (s *Server) initConfig() {
	cfg := config.NewConfig(s.log, s.OnConfigChange)
	s.OnConfigChange(cfg)
}

// OnConfigChange execute operation required on config change
func (s *Server) OnConfigChange(newCfg *config.AppConfig) {
	s.mu.Lock()
	defer func() { s.mu.Unlock() }()
	lvl := logger.ToLogLevel(newCfg.LogLevel)
	_ = s.log.SetLevel(lvl)
	s.config = newCfg
	s.initServices(context.Background())

	s.config.PrintConfig(s.log)
}

// Start starts web server main event loop
func (s *Server) Start() {
	// Determine API listen address/port
	serverListenAddr := fmt.Sprintf("0.0.0.0:%d", s.config.Port)
	// Start server
	go func() {
		if err := s.Echo.Start(serverListenAddr); err != nil {
			s.log.WithError(err).
				Fatal("Fatal error in API server")
		}
	}()
	logrus.Info(fmt.Sprintf("Service start listening on %s", serverListenAddr))
	// Wait for interrupt signal to gracefully shutting down the web server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancelShutdown()
	if err := s.Echo.Shutdown(ctx); err != nil {
		s.log.WithError(err).
			Fatal("Error shutting down API server")
	}
}
