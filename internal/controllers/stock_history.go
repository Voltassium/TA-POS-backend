package controllers

import (
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/services"
	"backend-ta/pkg/errors"
	internalHTTP "backend-ta/pkg/http"
	"backend-ta/pkg/http/server/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type StockHistoryController struct {
	service services.StockHistoryService
}

func NewStockHistoryController(service services.StockHistoryService) StockHistoryController {
	return StockHistoryController{service: service}
}

func (ctl *StockHistoryController) List(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	var productID string
	if err == nil {
		productID = id
	}

	var req requests.ListStockHistory
	if err := internalHTTP.BindData(ctx, &req); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	if productID != "" {
		req.ProductID = productID
	}

	res, err := ctl.service.List(ctx.Request.Context(), req)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}
