package httpServer

import (
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
)

func addRoutes(r *gin.Engine, env *Environment) {
	r.GET("/debug/vars", expvar.Handler())
	r.GET("/api/v0/views", func(c *gin.Context) {
		HandleViewsIndex(env, c)
	})

	// add dynamic routes
	for _, v := range env.Views {
		view := v
		r.GET("/"+view.Name(), func(c *gin.Context) {
			handleViewIndex(view, c)
		})

		for _, c := range view.Cameras() {
			camera := c

			cameraClient := env.CameraClientPoolInstance.GetClient(camera)
			if cameraClient == nil {
				continue
			}

			r.GET("/api/v0/images/"+view.Name()+"/"+camera+".jpg", func(c *gin.Context) {
				handleCameraImage(cameraClient, view, c)
			})
		}
	}
}
