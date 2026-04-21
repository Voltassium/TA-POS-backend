package http_response

import (
	internal_err "backend-ta/pkg/errors"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func SendSuccess(c *gin.Context, status int, message string, data any) {
	var msg string
	if message != "" {
		msg = message
	} else {
		msg = http.StatusText(status)
	}

	response := Response{
		Status:  status,
		Message: msg,
		Data:    data,
	}
	c.JSON(status, response)
}

func SendError(c *gin.Context, err error) {
	var appErr internal_err.AppError
	if errors.As(err, &appErr) {
		err = c.Error(appErr)
	} else {
		err = c.Error(err)
	}

	if err != nil {
		c.Abort()
	}
}
