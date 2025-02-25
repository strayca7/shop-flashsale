package adminhttp

import (
	"github.com/gin-gonic/gin"
)

func TestAdmin(r *gin.Engine) {
	r.GET("/admin", Ping)
}
