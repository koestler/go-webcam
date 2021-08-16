package httpServer

import (
	"github.com/gin-gonic/gin"
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

func getDimensions(view *config.ViewConfig, c *gin.Context) (dim Dimension) {
	dim.width = view.ResolutionMaxWidth()
	dim.height = view.ResolutionMaxHeight()

	if width := c.Query("width"); len(width) > 0 {
		if width, err := strconv.Atoi(width); err == nil {
			dim.width = min(dim.width, width)
		}
	}

	if height := c.Query("height"); len(height) > 0 {
		if height, err := strconv.Atoi(height); err == nil {
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
	c *gin.Context,
) {
	// fetch image
	cameraImage := cameraClient.GetResizedImage(view.RefreshInterval(), getDimensions(view, c))

	// handle camera fetching errors
	if cameraImage.Err() != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": cameraImage.Err()})
	}

	c.Header("Content-Type", "image/jpeg")
	c.Header("Access-Control-Allow-Origin", "*")

	c.Writer.Write(cameraImage.Img())
}
