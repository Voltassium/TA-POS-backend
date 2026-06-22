package routes

import (
	"backend-ta/app/constants"
	"backend-ta/app/controllers"
	"backend-ta/app/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerStatistics(router *gin.RouterGroup) {
	statsCtl := controllers.NewStatisticsController(services.ServicePool.StatisticsService)

	stats := router.Group("/statistics")
	{
		stats.Use(middlewares.RoleHandler(constants.UserRoleSuperadmin, constants.UserRoleOwner))
		stats.GET("/dashboard", statsCtl.GetDashboardData)
	}
}
