package httpServer

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
	"net/http"
)

func handleCameraImage(camera string, env *Environment, w http.ResponseWriter, r *http.Request) Error {
	// get camera client
	cameraClient := env.CameraClientPoolInstance.GetClient(camera)
	if cameraClient == nil {
		return StatusError{500, errors.New("camera client not started")}
	}

	// fetch image
	rawImage, err := cameraClient.GetRawImage()
	if err != nil {
		return StatusError{500, err}
	}

	// set headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", fmt.Sprintf(
		"public, max-age=%d",
		int(math.Floor(cameraClient.Config().RefreshInterval().Seconds()))),
	)

	// write image
	w.Write(rawImage)
	return nil
}
