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

type OrderItemController struct {
	orderService services.OrderService
}

func NewOrderItemController(orderSrv services.OrderService) OrderItemController {
	return OrderItemController{orderService: orderSrv}
}

func (ctl *OrderItemController) AddItem(ctx *gin.Context) {
	orderID, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.AddOrderItem
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.orderService.AddItem(ctx, orderID, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", res)
}

func (ctl *OrderItemController) RemoveItem(ctx *gin.Context) {
	orderID, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	itemID, err := internalHTTP.BindParams[int64](ctx, "item_id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.orderService.RemoveItem(ctx, orderID, itemID)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", res)
}
