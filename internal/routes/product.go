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
	stockHistoryCtl := controllers.NewStockHistoryController(services.ServicePool.StockHistoryService)

	product := router.Group("/products")
	{
		// Staff and Admin can view products and stock histories
		product.GET("", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), productCtl.List)
		product.GET(":id", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), productCtl.Get)
		product.GET(":id/stock-histories", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), stockHistoryCtl.List)

		// Admin only
		adminProduct := product.Group("")
		adminProduct.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		{
			adminProduct.POST("", productCtl.Create)
			adminProduct.PUT(":id", productCtl.Update)
			adminProduct.DELETE(":id", productCtl.Delete)
		}
	}

	stockHistories := router.Group("/stock-histories")
	{
		stockHistories.GET("", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), stockHistoryCtl.List)
	}
}
