package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type configResponse struct {
	ProjectTitle string
	Views        []viewResponse
}

type viewResponse struct {
	Name              string
	Title             string
	Cameras           []string
	RefreshIntervalMs int64
	Autoplay          bool
	Authentication    bool
}

// handleConfig godoc
// @Summary Show a account
// @Description get string by ID
// @Produce  json
// @Router /config [get]
func handleConfig(env *Environment, c *gin.Context) {
	response := configResponse{
		ProjectTitle: env.ProjectTitle,
		Views:        make([]viewResponse, 0),
	}

	for _, v := range env.Views {
		response.Views = append(response.Views, viewResponse{
			Name:              v.Name(),
			Title:             v.Title(),
			Cameras:           v.Cameras(),
			RefreshIntervalMs: v.RefreshInterval().Milliseconds(),
			Autoplay:          v.Autoplay(),
			Authentication:    false,
		})
	}

	c.JSON(http.StatusOK, response)
}
