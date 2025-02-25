package userhttp

import (
	"github.com/gin-gonic/gin"
)

func TestUser(r *gin.Engine) {
	r.GET("/user", Ping)
}