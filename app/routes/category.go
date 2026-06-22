package routes

import (
	"backend-ta/app/constants"
	"backend-ta/app/controllers"
	"backend-ta/app/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerCategory(router *gin.RouterGroup) {
	categoryCtl := controllers.NewCategoryController(services.ServicePool.CategoryService)

	category := router.Group("/categories")
	{
		// Staff and Admin can view categories
		category.GET("", middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner, constants.UserRoleChef, constants.UserRoleStaff), categoryCtl.List)
		category.GET(":id", middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner, constants.UserRoleChef, constants.UserRoleStaff), categoryCtl.Get)

		// Admin only
		adminCategory := category.Group("")
		adminCategory.Use(middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner))
		{
			adminCategory.POST("", categoryCtl.Create)
			adminCategory.PUT(":id", categoryCtl.Update)
			adminCategory.DELETE(":id", categoryCtl.Delete)
		}
	}
}
