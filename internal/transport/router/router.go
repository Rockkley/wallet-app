package router

import (
	"github.com/gin-gonic/gin"
	"wallet-app/internal/transport/handler"
)

func SetupRoutes(app *gin.Engine, handler *handler.WalletHandler) {
	api := app.Group("/api/v1")
	{
		api.POST("/wallet", handler.HandleTransaction)
		api.GET("/wallets/:id", handler.GetBalance)
	}
}
