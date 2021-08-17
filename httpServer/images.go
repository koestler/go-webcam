package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"log"
	"net/http"
	"strconv"
)

// setupImages godoc
// @Summary Outputs camera images.
// @Description Fetches the images from the camera (or from a cache), scales it to the requested resolution
// @Description and then returns it.
// @ID images
// @Param view path string true "View Name as provided by the config endpoint"
// @Param cameraName path string true "Camera Name as provided in Cameras array of the config endpoint"
// @Param width query int false "Downscale image to this width"
// @Param height query int false "Downscale image to this height"
// @Produce jpeg
// @Success 200
// @Failure 500 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /images/{view}/{cameraName}.jpg [get]
func setupImages(r *gin.RouterGroup, env *Environment) {
	// add dynamic routes
	for _, v := range env.Views {
		view := v
		for _, c := range view.Cameras() {
			camera := c

			client := env.CameraClientPoolInstance.GetClient(camera)
			if client == nil {
				continue
			}

			relativePath := "images/" + view.Name() + "/" + camera + ".jpg"
			r.GET(relativePath, func(c *gin.Context) {
				handleCameraImage(client, view, c)
			})
			log.Printf("httpServer: %s%s -> serve image", r.BasePath(), relativePath)
		}
	}
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
		NewErrorResponse(c, http.StatusServiceUnavailable, cameraImage.Err())
	}

	c.Header("Content-Type", "image/jpeg")

	c.Writer.Write(cameraImage.Img())
}

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
