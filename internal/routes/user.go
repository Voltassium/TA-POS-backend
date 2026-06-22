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
		user.GET("profile", middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner, constants.UserRoleChef, constants.UserRoleStaff), userCtl.GetProfile)

		adminUser := user.Group("")
		adminUser.Use(middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner))
		{
			adminUser.POST("", userCtl.CreateUserByAdmin)
			adminUser.GET("", userCtl.ListUser)
			adminUser.GET(":id", userCtl.Get)
			adminUser.PUT(":id", userCtl.Update)
			adminUser.DELETE(":id", userCtl.Delete)
		}
	}
}
