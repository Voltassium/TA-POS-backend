package controllers

import (
	"backend-ta/app/dto/requests"
	"backend-ta/app/services"
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
	orderID, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	itemID, err := internalHTTP.BindParams[string](ctx, "item_id")
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
