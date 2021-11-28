package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

type loginRequest struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

// setupLogin godoc
// @Summary Login endpoint
// @Description Creates a new JWT token used for authentication if a valud user / password is given.
// @ID login
// @Accept json
// @Produce json
// @Param request body loginRequest true "user info"
// @Success 200 {object} loginResponse
// @Failure 422 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func setupLogin(r *gin.RouterGroup, env *Environment) {
	r.POST("login", func(c *gin.Context) {
		var req loginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			jsonErrorResponse(c, http.StatusUnprocessableEntity, errors.New("Invalid json body provided"))
			return
		}

		if req.Password != "correct" {
			jsonErrorResponse(c, http.StatusUnauthorized, errors.New("Invalid credentials"))
			return
		}

		tokenStr, err := createJwtToken(env.Auth, req.User)
		if err != nil {
			jsonErrorResponse(c, http.StatusInternalServerError, errors.New("Cannot create token"))
			return
		}

		c.JSON(http.StatusOK, loginResponse{Token: tokenStr})
	})
	log.Printf("httpServer: %slogin -> serve login", r.BasePath())
}
