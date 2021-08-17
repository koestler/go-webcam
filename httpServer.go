package main

import (
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"github.com/koestler/go-webcam/httpServer"
	"log"
)

func runHttpServer(cfg *config.Config, cameraClientPoolInstance *cameraClient.ClientPool) *httpServer.HttpServer {
	httpServerCfg := cfg.HttpServer()
	if !httpServerCfg.Enabled() {
		return nil
	}

	if cfg.LogWorkerStart() {
		log.Printf("httpServer: start: bind=%s, port=%d", httpServerCfg.Bind(), httpServerCfg.Port())
	}

	return httpServer.Run(
		cfg.HttpServer(),
		&httpServer.Environment{
			Views:                    cfg.Views(),
			CameraClientPoolInstance: cameraClientPoolInstance,
		},
	)
}
