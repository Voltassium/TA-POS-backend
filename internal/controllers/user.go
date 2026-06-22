package controllers

import (
	"backend-ta/internal/dto/requests"
	"backend-ta/internal/services"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/errors"
	internalHTTP "backend-ta/pkg/http"
	"backend-ta/pkg/http/server/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(bankSrv services.UserService) UserController {
	return UserController{
		UserService: bankSrv,
	}
}

func (ctl *UserController) CreateUser(ctx *gin.Context) {
	var agent requests.CreateUser
	if err := internalHTTP.BindData(ctx, &agent); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.UserService.Register(ctx, agent)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "", nil)
}

func (ctl *UserController) CreateUserByAdmin(ctx *gin.Context) {
	var payload requests.CreateUserByAdmin
	if err := internalHTTP.BindData(ctx, &payload); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := ctl.UserService.RegisterByAdmin(ctx, payload)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusCreated, "Akun berhasil didaftarkan", nil)
}

func (ctl *UserController) ListUser(ctx *gin.Context) {
	var users requests.ListUser
	if err := internalHTTP.BindData(ctx, &users); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.UserService.GetList(ctx, users)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get list data", res)
}

func (ctl *UserController) Update(ctx *gin.Context) {
	var vendor requests.UpdateUser

	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	if err = internalHTTP.BindData(ctx, &vendor); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.UserService.Update(ctx, id, vendor)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *UserController) Delete(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err = ctl.UserService.DeleteSrv(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "", nil)
}

func (ctl *UserController) Get(ctx *gin.Context) {
	id, err := internalHTTP.BindParams[string](ctx, "id")
	if err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.UserService.Detail(ctx, id)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}

func (ctl *UserController) GetProfile(ctx *gin.Context) {
	res, err := ctl.UserService.Detail(ctx, authentication.GetUserDataFromToken(ctx).UserID)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	http_response.SendSuccess(ctx, http.StatusOK, "Success get data", res)
}
