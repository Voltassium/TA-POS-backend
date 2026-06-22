package routes

import (
	"backend-ta/app/controllers"
	"backend-ta/app/services"

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


