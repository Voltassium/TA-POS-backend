package routes

import (
	"backend-ta/internal/constants"
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerStatistics(router *gin.RouterGroup) {
	statsCtl := controllers.NewStatisticsController(services.ServicePool.StatisticsService)

	stats := router.Group("/statistics")
	{
		stats.Use(middlewares.RoleHandler(constants.UserRoleAdmin))
		stats.GET("/dashboard", statsCtl.GetDashboardData)
	}
}
