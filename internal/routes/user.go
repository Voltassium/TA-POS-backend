package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerUser(router *gin.RouterGroup) {
	userCtl := controllers.NewUserController(services.ServicePool.UserService)

	user := router.Group("/users")
	{
		user.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		user.GET("", userCtl.ListUser)
		user.GET("profile", userCtl.GetProfile)
		user.GET(":id", userCtl.Get)
		user.PUT(":id", userCtl.Update)
		user.DELETE(":id", userCtl.Delete)
	}
}
