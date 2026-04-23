package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerProduct(router *gin.RouterGroup) {
	productCtl := controllers.NewProductController(services.ServicePool.ProductService)

	product := router.Group("/products")
	{
		product.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		product.POST("", productCtl.Create)
		product.PUT(":id", productCtl.Update)
		product.DELETE(":id", productCtl.Delete)

		stockHistoryCtl := controllers.NewStockHistoryController(services.ServicePool.StockHistoryService)
		product.GET(":id/stock-histories", stockHistoryCtl.List)
	}
}

func registerProductPublic(router *gin.RouterGroup) {
	productCtl := controllers.NewProductController(services.ServicePool.ProductService)

	product := router.Group("/products")
	{
		product.GET("", productCtl.List)
		product.GET(":id", productCtl.Get)
	}
}
