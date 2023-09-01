// Package main application entrypoint
package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/shuvava/go-logging/logger"

	"github.com/shuvava/ota-tuf-server/internal/app"
	"github.com/shuvava/ota-tuf-server/pkg/version"
)

func main() {
	log := logger.NewLogrusLogger(logrus.InfoLevel)
	log.Info(fmt.Sprintf("Starting %s/%s", version.AppName, version.Version))
	log.Info(fmt.Sprintf("	Build date: %s", version.BuildDate))
	log.Info(fmt.Sprintf("	Commit hash: %s", version.CommitHash))

	server := app.NewServer(log)
	server.Start()
}
