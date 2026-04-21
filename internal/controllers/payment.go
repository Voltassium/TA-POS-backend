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

type PaymentController struct {
	paymentService services.PaymentService
}

func NewPaymentController(paymentSrv services.PaymentService) PaymentController {
	return PaymentController{paymentService: paymentSrv}
}

func (ctl *PaymentController) Create(ctx *gin.Context) {
	var payload requests.CreatePayment
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.paymentService.Process(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", res)
}

func (ctl *PaymentController) GetByOrder(ctx *gin.Context) {
	orderID, err := internalHTTP.BindParams[int64](ctx, "order_id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.paymentService.GetByOrder(ctx, orderID)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}
