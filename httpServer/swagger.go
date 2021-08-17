package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/docs"
	"github.com/swaggo/files"       // swagger embed files
	"github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title go-webcam API v0
// @version 0.0
// @description This server exposes webcam images to the public allowing for resizing, caching and authentication.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://github.com/koestler/go-webcam/blob/main/LICENSE

// @host http://localhost:8043/
// @BasePath /api/v0

func setupSwaggerDocs(r *gin.Engine, config Config) {
	docs.SwaggerInfo.Title = "Swagger Example API"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
