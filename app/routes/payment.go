package routes

import (
	"backend-ta/app/constants"
	"backend-ta/app/controllers"
	"backend-ta/app/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerPayment(router *gin.RouterGroup) {
	paymentCtl := controllers.NewPaymentController(services.ServicePool.PaymentService)

	payment := router.Group("/payments")
	payment.Use(middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner, constants.UserRoleStaff))
	{
		payment.POST("", paymentCtl.Create)
		payment.GET(":order_id", paymentCtl.GetByOrder)
	}
}
