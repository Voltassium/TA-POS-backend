package routes

import (
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"

	"github.com/gin-gonic/gin"
)

func registerStore(router *gin.RouterGroup) {
	storeCtl := controllers.NewStoreController(services.ServicePool.StoreService)

	storeRoutes := router.Group("/stores")
	{
		storeRoutes.POST("", storeCtl.Create)
		storeRoutes.GET("", storeCtl.List)
		storeRoutes.GET("/:id", storeCtl.Get)
		storeRoutes.PUT("/:id", storeCtl.Update)
		storeRoutes.DELETE("/:id", storeCtl.Delete)
	}
}
