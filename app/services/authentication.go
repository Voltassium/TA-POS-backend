package services

import (
	"backend-ta/app/constants"
	"backend-ta/app/dto/requests"
	"backend-ta/app/dto/response"
	"backend-ta/app/repository"
	"backend-ta/pkg/authentication"
	internal_err "backend-ta/pkg/errors"
	"context"
	"database/sql"
	"errors"
)

type AuthService interface {
	Login(ctx context.Context, payload requests.Login) (res response.LoginResponse, err error)
	Logout(ctx context.Context) (err error)
	RefreshToken(ctx context.Context, payload requests.RefreshToken) (res response.RefreshTokenResponse, err error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthSrv(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (a *authService) Login(ctx context.Context, payload requests.Login) (res response.LoginResponse, err error) {
	user, err := a.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
		}
		return res, err
	}

	isValid, err := authentication.VerifyPassword(payload.Password, user.Password)
	if err != nil {
		return res, err
	}

	if !isValid {
		return res, internal_err.AuthError(constants.AuthPasswordInvalidOrEmailNotFound)
	}

	tokenPayload := requests.ToTokenPayload(user)
	pair, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, false)
	if err != nil {
		return res, err
	}

	return response.LoginResponse{
		AccessToken:  pair.AccessToken,
		RefreshToken: pair.RefreshToken,
		Role:         string(user.Role),
	}, nil
}

func (a *authService) Logout(ctx context.Context) error {
	_ = authentication.GetUserDataFromToken(ctx)
	return nil
}

func (a *authService) RefreshToken(ctx context.Context, payload requests.RefreshToken) (res response.RefreshTokenResponse, err error) {
	claims, err := authentication.JWTAuth.VerifyRefreshToken(payload.RefreshToken)
	if err != nil {
		return res, err
	}
	if claims == nil {
		return res, internal_err.AuthError(constants.AuthInvalidToken)
	}

	user, err := a.userRepo.GetUser(ctx, claims.UserID)
	if err != nil {
		return res, err
	}

	tokenPayload := requests.ToTokenPayload(user)
	claimsRefresh, err := authentication.JWTAuth.GenerateTokenPair(tokenPayload, true)
	if err != nil {
		return res, err
	}

	return response.RefreshTokenResponse{
		AccessToken: claimsRefresh.AccessToken,
	}, nil
}
