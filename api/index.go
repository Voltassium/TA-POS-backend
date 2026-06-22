package api

import (
	"backend-ta/internal/routes"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/config"
	"backend-ta/pkg/database"
	"backend-ta/pkg/http/server/middlewares"
	"backend-ta/pkg/logger"
	"net/http"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func init() {
	cfg := config.LoadConfig()

	logger.NewZapLogger(cfg.Logger)
	database.InitDB(cfg.Database)

	authentication.NewJWTManager(authentication.JWTOptions{
		AccessSecret:       cfg.Authentication.AccessSecretKey,
		RefreshSecret:      cfg.Authentication.RefreshSecretKey,
		Issuer:             cfg.Authentication.Issuer,
		ExpiryAccessToken:  cfg.Authentication.AccessTokenExpiry,
		ExpiryRefreshToken: cfg.Authentication.RefreshTokenExpiry,
	})

	authentication.SetupKey(cfg.Authentication.EncryptKey)

	gin.SetMode(gin.ReleaseMode)

	r = gin.New()
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(middlewares.HandleCors())
	r.Use(middlewares.LoggerMiddleware())
	r.Use(middlewares.ErrorMiddleware())
	r.NoRoute(middlewares.NotFoundHandler)

	routes.RegisterV1(r)
}

func Handler(w http.ResponseWriter, req *http.Request) {
	r.ServeHTTP(w, req)
}
