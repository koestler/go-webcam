package httpServer

import (
	"fmt"
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"image/jpeg"
	"log"
	"math"
	"net/http"
	"strconv"
)

type Dimension struct {
	width  int
	height int
}

func (c Dimension) Width() int {
	return c.width
}

func (c Dimension) Height() int {
	return c.height
}

func getDimensions (view *config.ViewConfig, r *http.Request) (dim Dimension) {
	dim.width = view.ResolutionMaxWidth()
	dim.height = view.ResolutionMaxHeight()

	if list, ok := r.URL.Query()["width"]; ok {
		if width, err := strconv.Atoi(list[0]) ; err == nil {
			dim.width = min(dim.width, width)
		}
	}

	if list, ok := r.URL.Query()["height"]; ok {
		if height, err := strconv.Atoi(list[0]) ; err == nil {
			dim.height = min(dim.height, height)
		}
	}

	return
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func handleCameraImage(
	cameraClient *cameraClient.Client,
	view *config.ViewConfig,
	w http.ResponseWriter,
	r *http.Request,
) Error {
	// fetch image
	dim := getDimensions(view, r)

	image, err := cameraClient.GetResizedImage(dim)
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

	// todo: put image quality setting into view config
	err = jpeg.Encode(w, image, &jpeg.Options{Quality: 90})
	if err != nil {
		log.Printf("handleCameraImage failed: %v", err)
	}

	return nil
}
