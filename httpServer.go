package main

import (
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"github.com/koestler/go-webcam/hashStore"
	"github.com/koestler/go-webcam/httpServer"
	"log"
)

//go:generate swag init -g httpServer/swagger.go

func runHttpServer(cfg *config.Config, cameraClientPoolInstance *cameraClient.ClientPool) *httpServer.HttpServer {
	httpServerCfg := cfg.HttpServer()
	if !httpServerCfg.Enabled() {
		return nil
	}

	if cfg.LogWorkerStart() {
		log.Printf("httpServer: start: bind=%s, port=%d", httpServerCfg.Bind(), httpServerCfg.Port())
	}

	return httpServer.Run(
		// todo: refactor config and env into one object?
		httpServerConfig{
			cfg.HttpServer(),
			cfg.GetViewNames(),
			cfg.LogConfig(),
			cfg.LogDebug(),
		},
		&httpServer.Environment{
			ProjectTitle:             cfg.ProjectTitle(),
			Views:                    cfg.Views(),
			Auth:                     cfg.Auth(),
			CameraClientPoolInstance: cameraClientPoolInstance,
			HashStorage:              hashStore.Run(cfg.HttpServer()),
		},
	)
}

type httpServerConfig struct {
	config.HttpServerConfig
	viewNames []string
	logConfig bool
	logDebug  bool
}

func (c httpServerConfig) GetViewNames() []string {
	return c.viewNames
}

func (c httpServerConfig) LogConfig() bool {
	return c.logConfig
}

func (c httpServerConfig) LogDebug() bool {
	return c.logDebug
}
