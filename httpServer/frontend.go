package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http/httputil"
)

func setupFrontend(engine *gin.Engine, config Config) {
	frontendUrl := config.FrontendProxy()
	frontendPath := config.FrontendPath()

	if frontendUrl != nil {
		log.Printf("httpServer: setup frontend proxy to: %s", frontendUrl.String())
		engine.NoRoute(func(c *gin.Context) {
			proxy := httputil.NewSingleHostReverseProxy(frontendUrl)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
	} else {
		if len(frontendPath) > 0 {
			log.Printf("httpServer: serve frontend from local folder: %s", frontendPath)
		} else {
			log.Print("httpServer: no frontend configured")
		}
		engine.NoRoute(func(c *gin.Context) {
			NewErrorResponse(c, 404, errors.New("route not found"))
		})
	}
}
