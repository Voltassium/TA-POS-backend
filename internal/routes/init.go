package routes

import (
	"backend-ta/internal/repository"
	"backend-ta/internal/services"
	"backend-ta/pkg/database"
	"backend-ta/pkg/http/server/http_response"
	"backend-ta/pkg/http/server/middlewares"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterV1(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {

			http_response.SendSuccess(c, http.StatusOK, "pong", gin.H{
				"time": time.Now(),
				"ua":   c.Request.UserAgent(),
			})
		})

		//-----------Init Dependency Injection Pool---------//
		repository.Init(database.GetDB())
		services.Init()

		//-----------Public API-------------//
		registerAuth(v1)

		//-----------Private API-------------//
		private := v1.Group("")
		private.Use(middlewares.TokenAuthMiddleware())
		registerStore(private)
		registerUser(private)
		registerCategory(private)
		registerProduct(private)
		registerOrder(private)
		registerPayment(private)
		registerStatistics(private)
		registerPengeluaran(private)
	}

}
