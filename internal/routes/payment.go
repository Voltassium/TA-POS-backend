package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
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
