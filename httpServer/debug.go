package httpServer

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

var list = []string{
	"If-Modified-Since",
	"If-None-Match",
	"Cache-Control",
	"Pragma",
}

func debugHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		filteredHeaders := make(map[string][]string)
		for k, v := range c.Request.Header {
			if stringInSlice(k, list) {
				filteredHeaders[k] = v
			}
		}
		log.Printf("httpServer: %s, request headers: %v", c.Request.RequestURI, filteredHeaders)
		c.Next()
	}
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.ToUpper(b) == strings.ToUpper(a) {
			return true
		}
	}
	return false
}
