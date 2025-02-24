package router

import (
	"github.com/gin-gonic/gin"
	"shop-flashsale/internal/handler"
)

func ResRestRouter(r *gin.Engine){
	r.GET("/ping", handler.Ping)
}