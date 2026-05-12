package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerCategory(router *gin.RouterGroup) {
	categoryCtl := controllers.NewCategoryController(services.ServicePool.CategoryService)

	category := router.Group("/categories")
	{
		// Staff and Admin can view categories
		category.GET("", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), categoryCtl.List)
		category.GET(":id", middlewares.RoleHandler(constants.UserRoleAdmin, constants.UserRoleStaff), categoryCtl.Get)

		// Admin only
		adminCategory := category.Group("")
		adminCategory.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		{
			adminCategory.POST("", categoryCtl.Create)
			adminCategory.PUT(":id", categoryCtl.Update)
			adminCategory.DELETE(":id", categoryCtl.Delete)
		}
	}
}
