package httpServer

import (
	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine, env *Environment) {
	v0 := r.Group("/api/v0/")
	setupExpVar(v0)
	setupConfig(v0, env)
}

func addApiV0(r *gin.RouterGroup, env *Environment) {
	// add dynamic routes
	for _, v := range env.Views {
		view := v
		for _, c := range view.Cameras() {
			camera := c

			cameraClient := env.CameraClientPoolInstance.GetClient(camera)
			if cameraClient == nil {
				continue
			}

			r.GET("images/"+view.Name()+"/"+camera+".jpg", func(c *gin.Context) {
				handleCameraImage(cameraClient, view, c)
			})
		}
	}
}
