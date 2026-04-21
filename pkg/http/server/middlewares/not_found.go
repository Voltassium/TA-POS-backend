package middlewares

import (
	"backend-ta/pkg/errors"
	"backend-ta/pkg/http/server/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFoundHandler(c *gin.Context) {
	http_response.SendError(c, errors.NewDefaultError(http.StatusNotFound, "Not Found"))
}
