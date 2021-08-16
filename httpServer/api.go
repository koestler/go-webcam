package httpServer

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type view struct {
	Name        string
	Title       string
	Cameras     []string
	HasHtaccess bool
}

func HandleViewsIndex(env *Environment, c *gin.Context) {
	views := make([]view, 0)

	for _, v := range env.Views {
		views = append(views, view{
			Name:        v.Name(),
			Title:       v.Title(),
			Cameras:     v.Cameras(),
			HasHtaccess: false,
		})
	}

	c.JSON(http.StatusOK, views)
}
