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

type OrderController struct {
	orderService services.OrderService
}

func NewOrderController(orderSrv services.OrderService) OrderController {
	return OrderController{orderService: orderSrv}
}

func (ctl *OrderController) Create(ctx *gin.Context) {
	var payload requests.CreateOrder
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.orderService.Create(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", res)
}

func (ctl *OrderController) List(ctx *gin.Context) {
	var payload requests.ListOrder
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.orderService.List(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}

func (ctl *OrderController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.orderService.Detail(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *OrderController) UpdateStatus(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdateOrderStatus
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.orderService.UpdateStatus(ctx, id, payload); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *OrderController) Cancel(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.orderService.Cancel(ctx, id); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}
