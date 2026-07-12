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

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productSrv services.ProductService) ProductController {
	return ProductController{productService: productSrv}
}

func (ctl *ProductController) Create(ctx *gin.Context) {
	var payload requests.CreateProduct
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.productService.Create(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", res)
}

func (ctl *ProductController) List(ctx *gin.Context) {
	var payload requests.ListProduct
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.productService.List(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}

func (ctl *ProductController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.productService.Detail(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *ProductController) Update(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.UpdateProduct
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.productService.Update(ctx, id, payload); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *ProductController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.productService.Delete(ctx, id); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *ProductController) Restock(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	var payload requests.RestockProduct
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err := ctl.productService.Restock(ctx, id, payload); err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Berhasil melakukan kulakan stok", nil)
}

