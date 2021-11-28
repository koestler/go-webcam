package httpServer

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/koestler/go-webcam/config"
	"net/http"
	"time"
)

type jwtClaims struct {
	User string `json:"sub"`
	jwt.StandardClaims
}

func createJwtToken(config config.AuthConfig, user string) (tokenStr string, err error) {
	claims := &jwtClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.JwtValidityPeriod()).Unix(),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString(config.JwtSecret())
}

func authJwtMiddleware(env *Environment) gin.HandlerFunc {
	return func(c *gin.Context) {
		// extract jwt toke from authorization header if present
		tokenStr := c.GetHeader("Authorization")
		if len(tokenStr) < 1 {
			c.Next()
			return
		}

		// decode jwt token
		claims := &jwtClaims{}
		tkn, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return env.Auth.JwtSecret(), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// abort request if invalid token is given
		if !tkn.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// continue; if user is set this means a valid token is present
		c.Set("AuthUser", claims.User)
		c.Next()
	}
}
