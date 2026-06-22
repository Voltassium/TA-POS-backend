package routes

import (
	"backend-ta/app/constants"
	"backend-ta/app/controllers"
	"backend-ta/app/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerOrder(router *gin.RouterGroup) {
	orderCtl := controllers.NewOrderController(services.ServicePool.OrderService)
	orderItemCtl := controllers.NewOrderItemController(services.ServicePool.OrderService)
	kitchenCtl := controllers.NewKitchenController(services.ServicePool.KitchenService)

	order := router.Group("/orders")
	order.Use(middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner, constants.UserRoleChef, constants.UserRoleStaff))
	{
		order.GET("", orderCtl.List)
		order.GET(":id", orderCtl.Get)
		order.POST("", orderCtl.Create)
		order.PATCH(":id/status", orderCtl.UpdateStatus)
		order.DELETE(":id", orderCtl.Cancel)

		order.POST(":id/items", orderItemCtl.AddItem)
		order.DELETE(":id/items/:item_id", orderItemCtl.RemoveItem)
		order.PATCH(":id/items/:item_id/served", kitchenCtl.UpdateItemServedQty)
	}
}
