package main

import (
	"shop-flashsale/internal/admin/delivery/http"
	"shop-flashsale/internal/user/delivery/http"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	adminhttp.TestAdmin(r)
	userhttp.TestUser(r)
	r.Run(":8080")
}
