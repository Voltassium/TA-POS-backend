package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerOrder(router *gin.RouterGroup) {
	orderCtl := controllers.NewOrderController(services.ServicePool.OrderService)
	orderItemCtl := controllers.NewOrderItemController(services.ServicePool.OrderService)

	order := router.Group("/orders")
	order.Use(middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff))
	{
		order.GET("", orderCtl.List)
		order.GET(":id", orderCtl.Get)
		order.POST("", orderCtl.Create)
		order.PATCH(":id/status", orderCtl.UpdateStatus)
		order.DELETE(":id", orderCtl.Cancel)

		order.POST(":id/items", orderItemCtl.AddItem)
		order.DELETE(":id/items/:item_id", orderItemCtl.RemoveItem)
	}
}
