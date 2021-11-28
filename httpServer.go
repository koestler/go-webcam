package main

import (
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
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
		httpServerConfig{
			cfg.HttpServer(),
			cfg.GetViewNames(),
		},
		&httpServer.Environment{
			ProjectTitle:             cfg.ProjectTitle(),
			Views:                    cfg.Views(),
			Auth:                     cfg.Auth(),
			CameraClientPoolInstance: cameraClientPoolInstance,
		},
	)
}

type httpServerConfig struct {
	config.HttpServerConfig
	viewNames []string
}

func (c httpServerConfig) GetViewNames() []string {
	return c.viewNames
}
