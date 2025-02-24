package main

import (
	"github.com/gin-gonic/gin"
	"shop-flashsale/internal/router"
)
func main() {
	r := gin.Default()
	router.ResRestRouter(r)
	r.Run(":8080")
}