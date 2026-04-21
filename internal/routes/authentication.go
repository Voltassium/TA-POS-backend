package routes

import (
	"backend-ta/internal/controllers"
	"backend-ta/internal/services"
	"backend-ta/pkg/http/server/middlewares"

	"github.com/gin-gonic/gin"
)

func registerAuth(router *gin.RouterGroup) {
	authCtl := controllers.NewAuthController(services.ServicePool.AuthService)
	userCtl := controllers.NewUserController(services.ServicePool.UserService)

	auth := router.Group("/authentications")
	{
		auth.POST("register", userCtl.CreateUser)
		auth.POST("login", authCtl.Login)
		auth.POST("refresh-token", authCtl.RefreshToken)

		auth.Use(middlewares.TokenAuthMiddleware())
		auth.POST("logout", authCtl.Logout)
	}
}
