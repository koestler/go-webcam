package main

import (
	"github.com/koestler/go-webcam/config"
	"github.com/koestler/go-webcam/httpServer"
	"log"
)

func runHttpServer(cfg *config.Config) *httpServer.HttpServer {
	if !cfg.HttpServer.Enabled() {
		return nil
	}

	if cfg.LogWorkerStart {
		log.Printf("httpServer: start: bind=%s, port=%d", cfg.HttpServer.Bind(), cfg.HttpServer.Port())
	}

	return httpServer.Run(
		&cfg.HttpServer,
		&httpServer.Environment{
			Views: cfg.Views,
		},
	)
}
