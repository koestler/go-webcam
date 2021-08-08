package httpServer

import (
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
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

func getDimensions(view *config.ViewConfig, r *http.Request) (dim Dimension) {
	dim.width = view.ResolutionMaxWidth()
	dim.height = view.ResolutionMaxHeight()

	if list, ok := r.URL.Query()["width"]; ok {
		if width, err := strconv.Atoi(list[0]); err == nil {
			dim.width = min(dim.width, width)
		}
	}

	if list, ok := r.URL.Query()["height"]; ok {
		if height, err := strconv.Atoi(list[0]); err == nil {
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
	cameraImage := cameraClient.GetDelayedImage(view.RefreshInterval(), getDimensions(view, r))

	// handle camera fetching errors
	if cameraImage.Err() != nil {
		return StatusError{http.StatusServiceUnavailable, cameraImage.Err()}
	}

	// set headers
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.Write(cameraImage.Img())

	return nil
}
