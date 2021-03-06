package app

import (
	"context"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	cmnapi "github.com/shuvava/go-ota-svc-common/api"

	"github.com/shuvava/ota-tuf-server/internal/api"
	"github.com/shuvava/ota-tuf-server/pkg/version"
)

const (
	routeAPIVer1 = "/api/v1"
)

// initWebServer creates echo http server and set request handlers
func (s *Server) initWebServer() {
	// Initialize Echo, set error handler, add in middleware
	e := echo.New()

	e.HideBanner = true
	e.HidePort = true
	e.HTTPErrorHandler = cmnapi.NewErrorHandler().Handler
	e.Pre(middleware.RemoveTrailingSlash())
	// logger Middleware (https://echo.labstack.com/middleware/logger/)
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Gzip())
	// Server header
	e.Use(cmnapi.ServerHeader(version.AppName, version.Version))
	initHealthRoutes(s, e)
	v1Group := e.Group(routeAPIVer1, middleware.RequestID())
	initKeyRepoRoutes(s, v1Group)

	// Enable metrics middleware
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	s.Echo = e
}

func initKeyRepoRoutes(s *Server, group *echo.Group) {
	group.POST(api.PathCreateRoot, func(c echo.Context) error {
		return api.CreateRoot(c, s.svc.KeySvc)
	})
}

func initHealthRoutes(s *Server, e *echo.Echo) {
	// Define a separate root 'health' group without the logging middleware added (for healthz/readyz)
	healthGroup := e.Group("")
	healthGroup.GET(cmnapi.LivenessPath, cmnapi.HealthzHandler)
	healthGroup.GET(cmnapi.ReadinessPath, cmnapi.ReadyzHandler(
		func(ctx context.Context) cmnapi.HealthEntryStatus {
			resource := "repository"
			if err := s.svc.Db.Ping(ctx); err != nil {
				return cmnapi.HealthEntryStatus{
					Status:   cmnapi.StatusUnhealthy,
					Data:     err.Error(),
					Resource: resource,
				}
			}
			return cmnapi.HealthEntryStatus{
				Status:   cmnapi.StatusHealthy,
				Resource: resource,
			}
		},
	))
}
