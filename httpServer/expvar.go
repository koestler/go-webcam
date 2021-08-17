package httpServer

import (
	"github.com/gin-contrib/expvar"
	"github.com/gin-gonic/gin"
)

// setupExpVar godoc
// @Summary Provide operation counters.
// @Description Package expvar provides a standardized interface to public variables,
// @Description such as operation counters in servers.
// @Description It exposes these variables via HTTP at /debug/vars in JSON format.
// @ID expvar
// @Produce json
// @Success 200
// @Router /debug/vars [get]
func setupExpVar(r *gin.RouterGroup) {
	r.GET("debug/vars", expvar.Handler())
}
