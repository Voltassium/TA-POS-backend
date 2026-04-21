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
		category.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		category.POST("", categoryCtl.Create)
		category.PUT(":id", categoryCtl.Update)
		category.DELETE(":id", categoryCtl.Delete)
	}
}

func registerCategoryPublic(router *gin.RouterGroup) {
	categoryCtl := controllers.NewCategoryController(services.ServicePool.CategoryService)

	category := router.Group("/categories")
	{
		category.GET("", categoryCtl.List)
		category.GET(":id", categoryCtl.Get)
	}
}
