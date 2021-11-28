package httpServer

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Message string `json:"message" example:"status bad request"`
}

func jsonErrorResponse(c *gin.Context, status int, err error) {
	er := ErrorResponse{
		Message: err.Error(),
	}
	c.JSON(status, er)
}

func jsonGetResponse(c *gin.Context, obj interface{}) {
	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		c.Status(http.StatusInternalServerError)
	}
	etag := fmt.Sprintf("%x", md5.Sum(jsonBytes))

	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}

	c.Header("ETag", "W/"+etag)
	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonBytes)
}
