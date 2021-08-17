package httpServer

import (
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine, env *Environment) {
	v0 := r.Group("/api/v0/")
	addApiV0(v0, env)
}

func addApiV0(r *gin.RouterGroup, env *Environment) {
	r.GET("debug/vars", expvar.Handler())
	r.GET("config", func(c *gin.Context) {
		handleConfig(env, c)
	})

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
