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

type StoreController struct {
	storeService services.StoreService
}

func NewStoreController(storeSrv services.StoreService) StoreController {
	return StoreController{storeService: storeSrv}
}

func (ctl *StoreController) Create(ctx *gin.Context) {
	var payload requests.CreateStore

	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.storeService.CreateStore(ctx.Request.Context(), payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "Store created successfully", res)
}

func (ctl *StoreController) List(ctx *gin.Context) {
	page := 1
	pageSize := 10
	orderBy := "created_at"
	orderDir := "desc"

	type ListReq struct {
		Page     int    `form:"page,default=1"`
		PageSize int    `form:"page_size,default=10"`
		OrderBy  string `form:"order_by,default=created_at"`
		OrderDir string `form:"order_dir,default=desc"`
	}
	var req ListReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		page = req.Page
		pageSize = req.PageSize
		orderBy = req.OrderBy
		orderDir = req.OrderDir
	}

	res, total, err := ctl.storeService.ListStores(ctx.Request.Context(), page, pageSize, orderBy, orderDir)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Stores listed successfully", map[string]interface{}{
		"data": res,
		"meta": map[string]interface{}{
			"total":     total,
			"page":      page,
			"page_size": pageSize,
		},
	})
}

func (ctl *StoreController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.storeService.GetStore(ctx.Request.Context(), id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Store fetched successfully", res)
}

func (ctl *StoreController) Update(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdateStore
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.storeService.UpdateStore(ctx.Request.Context(), id, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Store updated successfully", res)
}

func (ctl *StoreController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[int64](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.storeService.DeleteStore(ctx.Request.Context(), id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Store deleted successfully", nil)
}
