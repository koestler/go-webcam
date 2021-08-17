package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http/httputil"
	"os"
	"path"
	"path/filepath"
)

func setupFrontend(engine *gin.Engine, config Config) {
	frontendUrl := config.FrontendProxy()

	if frontendUrl != nil {
		engine.NoRoute(func(c *gin.Context) {
			proxy := httputil.NewSingleHostReverseProxy(frontendUrl)
			proxy.ServeHTTP(c.Writer, c.Request)
		})
		log.Printf("httpServer: /* -> proxy %s*", frontendUrl)
	} else {
		frontendPath := path.Clean(config.FrontendPath())

		if len(frontendPath) > 0 {
			if frontendPathInfo, err := os.Lstat(frontendPath); err != nil {
				log.Printf("httpServer: given frontend path is not accessible: %s", err)
			} else if !frontendPathInfo.IsDir() {
				log.Printf("httpServer: given frontend path is not a directory")
			} else {
				err := filepath.Walk(frontendPath, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if info.IsDir() {
						return nil
					}

					route := path[len(frontendPath):]
					if route == "/index.html" {
						route = "/"
					}
					engine.StaticFile(route, path)
					log.Printf("httpServer: %s -> serve %s", route, path)

					return nil
				})

				if err != nil {
					log.Printf("httpServer: failed to serve from local folder: %s", err)
				}
			}
		} else {
			log.Print("httpServer: no frontend configured")
		}
		engine.NoRoute(func(c *gin.Context) {
			NewErrorResponse(c, 404, errors.New("route not found"))
		})
	}
}
