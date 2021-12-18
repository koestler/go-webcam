package httpServer

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/koestler/go-webcam/cameraClient"
	"log"
	"math"
	"math/rand"
	"net/http"
	"regexp"
)

// setupImagesByHash godoc
// @Summary Outputs camera images.
// @Description Returns an image of which the hash is known.
// @ID imagesByHash
// @Param hash path string true "The hash of the image properties."
// @Produce jpeg
// @Success 200
// @Failure 500 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /imagesByHash/{hash}.jpg [get]
// @Security ApiKeyAuth
func setupImagesByHash(r *gin.RouterGroup, env *Environment) {
	r.GET("imagesByHash/:filename", func(c *gin.Context) {
		filename := c.Param("filename")
		if !hashFileNameMatcher.MatchString(filename) {
			jsonErrorResponse(c, http.StatusNotFound, fmt.Errorf("invalid filename: '%s", filename))
			return
		}
		hash := filename[0:40]
		cp := env.HashStorage.Get(hash)
		if cp == nil {
			jsonErrorResponse(c, http.StatusNotFound, fmt.Errorf("unknown hash: '%s", hash))
			return
		}
		maxAge := math.Ceil(env.HashStorage.Config().HashTimeout().Seconds())
		c.Header("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
		c.Data(http.StatusOK, "image/jpeg", cp.JpgImg())
	})
	log.Printf("httpServer: %simagesByHash/<hash>.jpg -> serve imagesByHash", r.BasePath())
}

var hashFileNameMatcher = regexp.MustCompilePOSIX(`^[0-9a-f]{40}\.jpg$`)

func getImageByHashUrl(cp cameraClient.CameraPicture, env *Environment) string {
	hash := getHash(cp)
	env.HashStorage.Set(hash, cp)
	return fmt.Sprintf("/api/v0/imagesByHash/%s.jpg", hash)
}

func getHash(cp cameraClient.CameraPicture) string {
	dimKey := cameraClient.DimensionCacheKey(cameraClient.DimensionOfImage(cp.DecodedImg()))
	str := fmt.Sprintf("%s-%s-%s", randomPrefix, cp.Uuid(), dimKey)
	h := sha1.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

var randomPrefix = randomString(64)

func randomString(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	s := make([]rune, n)
	for i := range s {
		s[i] = letters[rand.Intn(len(letters))]
	}
	return string(s)
}
