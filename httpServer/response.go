package httpServer

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
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
	hash := md5.Sum(jsonBytes)
	etag := hex.EncodeToString(hash[:])
	if match := c.GetHeader("If-None-Match"); match != "" {
		if strings.Contains(match, etag) {
			c.AbortWithStatus(http.StatusNotModified)
			return
		}
	}

	c.Header("ETag", "W/"+etag)
	c.Data(http.StatusOK, "application/json; charset=utf-8", jsonBytes)
}

func setCacheControlPublic(c *gin.Context, maxAge time.Duration) {
	c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", int(maxAge.Seconds())))
}
