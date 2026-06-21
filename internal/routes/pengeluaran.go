package routes

import (
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"

	"github.com/gin-gonic/gin"
)

func registerPengeluaran(router *gin.RouterGroup) {
	ctl := controllers.NewPengeluaranController(services.ServicePool.PengeluaranService)

	routes := router.Group("/pengeluaran")
	{
		routes.POST("", ctl.Create)
		routes.GET("", ctl.List)
		routes.GET("/:id", ctl.Get)
		routes.PUT("/:id", ctl.Update)
		routes.DELETE("/:id", ctl.Delete)
	}
}


