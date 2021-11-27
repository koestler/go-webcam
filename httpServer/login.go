package httpServer

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// setupLogin godoc
// @Summary Login endpoint
// @Description Creates a new JWT token used for authentication if a valud user / password is given.
// @ID login
// @Produce json
// @Success 200 {object} loginResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func setupLogin(r *gin.RouterGroup, env *Environment) {
	r.POST("login", func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&loginRequest{}); err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Invalid json body provided")
			return
		}

		if req.Password != "correct" {
			c.JSON(http.StatusUnauthorized, "Invalid credentials")
			return
		}

		//token, err := createToken(env, req.Username)
		//if err != nil {
//			c.JSON(http.StatusUnprocessableEntity, err.Error())
//			return
		//}

		jsonResponse(c, loginResponse{			Token: "abcd"	})
	})
	log.Printf("httpServer: %slogin -> serve login", r.BasePath())
}
