package controllers

import (
	"backend-ta/app/constants"
	"backend-ta/app/dto/requests"
	"backend-ta/app/services"
	"backend-ta/pkg/errors"
	internalHTTP "backend-ta/pkg/http"
	"backend-ta/pkg/http/server/http_response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService services.AuthService
}

func NewAuthController(authSrv services.AuthService) AuthController {
	return AuthController{
		authService: authSrv,
	}
}

func (ctl *AuthController) Login(ctx *gin.Context) {
	var auth requests.Login
	if err := internalHTTP.BindData(ctx, &auth); err != nil {
		http_response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	res, err := ctl.authService.Login(ctx, auth)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)

	ctx.SetCookie("access_token", res.AccessToken, 86400, "/", "", true, true)
	ctx.SetCookie("refresh_token", res.RefreshToken, 604800, "/", "", true, true)

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthLoginSuccess, res)
}

func (ctl *AuthController) Logout(ctx *gin.Context) {
	err := ctl.authService.Logout(ctx)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)

	ctx.SetCookie("access_token", "", -1, "/", "", true, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "", true, true)

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthLogoutSuccess, nil)
}

func (ctl *AuthController) RefreshToken(ctx *gin.Context) {
	token, err := ctx.Cookie("refresh_token")
	if err != nil || token == "" {
		http_response.SendError(ctx, errors.AuthError("missing refresh token"))
		return
	}

	auth := requests.RefreshToken{RefreshToken: token}

	res, err := ctl.authService.RefreshToken(ctx, auth)
	if err != nil {
		http_response.SendError(ctx, err)
		return
	}

	ctx.SetSameSite(http.SameSiteNoneMode)

	ctx.SetCookie("access_token", res.AccessToken, 86400, "/", "", true, true)

	http_response.SendSuccess(ctx, http.StatusOK, constants.AuthRefreshTokenSuccess, res)
}
