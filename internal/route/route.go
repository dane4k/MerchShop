package route

import (
	"MerchShop/app"
	"github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine, appObj *app.App) {
	api := router.Group("/api")
	api.Use(appObj.AuthMiddleware.Handle())
	{
		api.GET("/info", appObj.InfoHandler.GetInfo)
		api.POST("/sendCoin", appObj.TransactionHandler.SendCoins)
		api.GET("/buy/:item", appObj.InventoryHandler.BuyItem)
	}

	router.POST("/api/auth", appObj.AuthHandler.Auth)
}
