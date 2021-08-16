package httpServer

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/config"
	"net/http"
)

func handleViewIndex(view *config.ViewConfig, c *gin.Context) {
	ret := fmt.Sprintf("<h1>%s</h1>", view.Name())
	ret += "<ul>"
	for _, camera := range view.Cameras() {
		ret += fmt.Sprintf("<li><img src=\"/api/v0/images/%s/%s.jpg?width=400\" /></li>", view.Name(), camera)
	}

	ret += "</ul>"

	c.Data(http.StatusOK, "text/html; charset=utf-8",[]byte(ret))
}
