package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitializeRouter() *gin.Engine {
	router := gin.Default()
	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{"message": "page not found"}) })
	router.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "up and running...") })

	transactions := router.Group("/transactions")
	transactions.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, gin.H{"status": "success"}) })
	return router
}
