package httpServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title go-webcam API v0
// @version 0.0
// @description This server fetches images from ip cameras and serves them over http.
//   It allows for resizing, caching and authentication.

// @license.name MIT
// @license.url https://github.com/koestler/go-webcam/blob/main/LICENSE

// @BasePath /api/v0
func setupSwaggerDocs(r *gin.Engine, config Config) {
	docs.SwaggerInfo.Host = fmt.Sprintf("127.0.0.1:%d", config.Port())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
