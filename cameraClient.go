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

	for _, cameraClientConfig := range cfg.Cameras {
		if cfg.LogWorkerStart {
			log.Printf(
				"cameraClient[%s]: start: address='%s'",
				cameraClientConfig.Name(),
				cameraClientConfig.Address(),
			)
		}

		if client, err := cameraClient.RunClient(cameraClientConfig); err != nil {
			log.Printf("cameraClient[%s]: start failed: %s", cameraClientConfig.Name(), err)
		} else {
			cameraClientPoolInstance.AddClient(client)
			countStarted += 1
			if cfg.LogWorkerStart {
				log.Printf(
					"cameraClient[%s]: started",
					cameraClientConfig.Name(),
				)
			}
		}
	}

	if countStarted < 1 {
		initiateShutdown <- errors.New("no cameraClient was started")
	}

	return cameraClientPoolInstance
}
