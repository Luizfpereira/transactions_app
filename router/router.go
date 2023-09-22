package router

import (
	"net/http"
	"transactions_app/handlers"

	"github.com/gin-gonic/gin"
)

func InitializeRouter(transactionHandler *handlers.TransactionHandler) *gin.Engine {
	router := gin.Default()
	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{"message": "page not found"}) })
	router.GET("/", func(ctx *gin.Context) { ctx.JSON(http.StatusOK, "up and running...") })

	transactions := router.Group("/transactions")
	transactions.POST("/create", transactionHandler.Create)
	return router
}
