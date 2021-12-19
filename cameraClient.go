package main

import (
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"github.com/pkg/errors"
	"log"
)

func runCameraClient(
	cfg *config.Config,
	initiateShutdown chan<- error,
) *cameraClient.ClientPool {
	cameraClientPoolInstance := cameraClient.RunPool()

	countStarted := 0

	for _, camera := range cfg.Cameras() {
		if cfg.LogWorkerStart() {
			log.Printf(
				"cameraClient[%s]: start: address='%s'",
				camera.Name(),
				camera.Address(),
			)
		}

		cameraConfig := cameraClientConfig{
			CameraConfig: *camera,
			logDebug:     cfg.LogDebug(),
		}

		if client, err := cameraClient.RunClient(&cameraConfig); err != nil {
			log.Printf("cameraClient[%s]: start failed: %s", camera.Name(), err)
		} else {
			cameraClientPoolInstance.AddClient(client)
			countStarted += 1
			if cfg.LogWorkerStart() {
				log.Printf(
					"cameraClient[%s]: started",
					camera.Name(),
				)
			}
		}
	}

	if countStarted < 1 {
		initiateShutdown <- errors.New("no cameraClient was started")
	}

	return cameraClientPoolInstance
}

type cameraClientConfig struct {
	config.CameraConfig
	logDebug bool
}

func (cc *cameraClientConfig) LogDebug() bool {
	return cc.logDebug
}