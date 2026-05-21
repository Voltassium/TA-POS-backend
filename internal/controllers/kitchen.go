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

type KitchenController struct {
	kitchenService services.KitchenService
}

func NewKitchenController(kitchenSrv services.KitchenService) KitchenController {
	return KitchenController{kitchenService: kitchenSrv}
}

func (ctl *KitchenController) UpdateItemServedQty(ctx *gin.Context) {
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

	var payload requests.UpdateServedQty
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.kitchenService.UpdateItemServedQty(ctx, orderID, itemID, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", res)
}
