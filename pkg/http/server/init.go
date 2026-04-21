package server

import (
	"backend-ta/pkg/config"
	"backend-ta/pkg/http/server/middlewares"
	"backend-ta/pkg/logger"
	"backend-ta/pkg/validation"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type HTTPServer struct {
	server *http.Server
}
type RegisterRoute func(*gin.Engine)

func Init(config config.Application, routes ...RegisterRoute) *HTTPServer {
	router := gin.New()
	router.Use(gzip.Gzip(gzip.DefaultCompression))
	router.Use(middlewares.HandleCors())
	router.Use(middlewares.LoggerMiddleware())
	router.Use(middlewares.ErrorMiddleware())
	router.NoRoute(middlewares.NotFoundHandler)

	//init router
	for _, route := range routes {
		route(router)
	}

	validation.InitGinValidator()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Error("Failed to start server:", zap.Error(err))
		}
	}()

	return &HTTPServer{
		server: srv,
	}
}

func (h *HTTPServer) GracefulShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	logger.Log.Info("Shutting down gracefully...")

	ctxShutDown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := h.server.Shutdown(ctxShutDown); err != nil {
		logger.Log.Error("Server forced to shutdown:", zap.Error(err))
	} else {
		logger.Log.Info("Server gracefully stopped")
	}
}
