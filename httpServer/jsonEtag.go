package httpServer

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func jsonResponse(c *gin.Context, obj interface{}) {
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
