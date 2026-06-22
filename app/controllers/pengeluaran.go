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

type PengeluaranController struct {
	service services.PengeluaranService
}

func NewPengeluaranController(srv services.PengeluaranService) PengeluaranController {
	return PengeluaranController{service: srv}
}

func (ctl *PengeluaranController) Create(ctx *gin.Context) {
	var payload requests.CreatePengeluaran
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.service.Create(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", res)
}

func (ctl *PengeluaranController) List(ctx *gin.Context) {
	var payload requests.ListPengeluaran
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.service.List(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}

func (ctl *PengeluaranController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.service.Detail(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *PengeluaranController) Update(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdatePengeluaran
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.service.Update(ctx, id, payload); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *PengeluaranController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.service.Delete(ctx, id); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}
