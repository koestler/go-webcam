package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/cameraClient"
	"github.com/koestler/go-webcam/config"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strconv"
	"time"
)

// setupImages godoc
// @Summary Outputs camera images.
// @Description Fetches the images from the camera (or from a cache), scales it to the requested resolution
// @Description and then returns it.
// @ID images
// @Param viewName path string true "View Name as provided by the config endpoint"
// @Param cameraName path string true "Camera Name as provided in Cameras array of the config endpoint"
// @Param width query int false "Downscale image to this width"
// @Param height query int false "Downscale image to this height"
// @Produce jpeg
// @Success 200
// @Success 307
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /images/{viewName}/{cameraName}.jpg [get]
// @Security ApiKeyAuth
func setupImages(r *gin.RouterGroup, env *Environment) {
	// add dynamic routes
	for _, v := range env.Views {
		view := v
		for _, c := range view.CameraNames() {
			camera := c

			client := env.CameraClientPoolInstance.GetClient(camera)
			if client == nil {
				continue
			}

			relativePath := "images/" + view.Name() + "/" + camera + ".jpg"
			r.GET(relativePath, func(c *gin.Context) {
				handleCameraImage(client, view, c, env)
			})
			if env.Config.LogConfig() {
				log.Printf("httpServer: %s%s -> serve image", r.BasePath(), relativePath)
			}
		}
	}
}

func handleCameraImage(
	cameraClient *cameraClient.Client,
	view *config.ViewConfig,
	c *gin.Context,
	env *Environment,
) {
	// check authorization
	if !isAuthenticated(view, c) {
		jsonErrorResponse(c, http.StatusForbidden, errors.New("User is not allowed here"))
		return
	}

	// fetch image
	cameraPicture := cameraClient.GetResizedImage(
		view.RefreshInterval(), getDimensions(view, c),
		view.JpgQuality(),
	)

	// handle camera fetching errors
	if cameraPicture.Err() != nil {
		jsonErrorResponse(c, http.StatusServiceUnavailable, cameraPicture.Err())
		return
	}

	if view.IsPublic() {
		// expire proxy cache before image expires.
		// The accounts for the fact that some proxies count the maxAge when they fully received the body
		// but this header is computed when the head is sent.
		setCacheControlPublicProxy(c, time.Until(cameraPicture.Expires())-env.Config.ImageEarlyExpire())
	}
	c.Redirect(http.StatusTemporaryRedirect, getImageByHashUrl(cameraPicture, env))
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

func isAuthenticated(view *config.ViewConfig, c *gin.Context) bool {
	if view.IsPublic() {
		return true
	}

	user := c.GetString("AuthUser")
	if len(user) < 1 {
		return false
	}

	return view.IsAllowed(user)
}
