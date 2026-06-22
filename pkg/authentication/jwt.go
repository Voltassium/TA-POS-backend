package authentication

import (
	"backend-ta/app/constants"
	"backend-ta/app/dto/requests"
	"backend-ta/pkg/errors"
	"context"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	JWTAuth *JWTManager
	once    = &sync.Once{}
)

type (
	AccessTokenClaims struct {
		requests.UserAuth
		jwt.RegisteredClaims
	}
	RefreshTokenClaims struct {
		UserID string `json:"user_id"`
		jwt.RegisteredClaims
	}
	TokenPair struct {
		AccessToken  string
		RefreshToken string
	}

	JWTManager struct {
		accessSecret                          []byte
		refreshSecret                         []byte
		issuer                                string
		expiryAccessToken, expiryRefreshToken time.Duration
	}

	JWTOptions struct {
		AccessSecret, RefreshSecret, Issuer   string
		ExpiryAccessToken, ExpiryRefreshToken time.Duration
	}
)

const Token string = "token"

func NewJWTManager(options JWTOptions) {
	once.Do(func() {
		JWTAuth = &JWTManager{
			accessSecret:       []byte(options.AccessSecret),
			refreshSecret:      []byte(options.RefreshSecret),
			issuer:             options.Issuer,
			expiryAccessToken:  options.ExpiryAccessToken,
			expiryRefreshToken: options.ExpiryRefreshToken,
		}
	})
}

func (m *JWTManager) GenerateTokenPair(payload requests.UserAuth, isAccessTokenOnly bool) (*TokenPair, error) {
	accessToken, err := m.generateAccessToken(payload)
	if err != nil {
		return nil, err
	}

	var refreshToken string
	if !isAccessTokenOnly {
		refreshToken, err = m.generateRefreshToken(payload.UserID, payload.Email)
		if err != nil {
			return nil, err
		}
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (m *JWTManager) generateAccessToken(auth requests.UserAuth) (string, error) {
	expirationTime := time.Now().Add(m.expiryAccessToken)
	claims := &AccessTokenClaims{
		UserAuth: requests.UserAuth{
			UserID:  auth.UserID,
			StoreID: auth.StoreID,
			Email:   auth.Email,
			Role:    auth.Role,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
			Subject:   auth.Email,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.accessSecret)
}

func (m *JWTManager) generateRefreshToken(userID string, subject string) (string, error) {
	expirationTime := time.Now().Add(m.expiryRefreshToken)

	claims := &RefreshTokenClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    m.issuer,
			Subject:   subject,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(m.refreshSecret)
}

func (m *JWTManager) VerifyAccessToken(tokenString string) (*AccessTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.AuthError(AuthErrSigningMethod.Error())
		}
		return m.accessSecret, nil
	})

	if err != nil {
		return nil, errors.AuthError(err.Error())
	}

	claims, ok := token.Claims.(*AccessTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.AuthError(AuthErrInvalidToken.Error())
	}

	return claims, nil
}

func (m *JWTManager) VerifyRefreshToken(tokenString string) (*RefreshTokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.AuthError(AuthErrSigningMethod.Error())
		}
		return m.refreshSecret, nil
	})

	if err != nil {
		return nil, errors.AuthError(err.Error())
	}

	claims, ok := token.Claims.(*RefreshTokenClaims)
	if !ok || !token.Valid {
		return nil, errors.AuthError(AuthErrInvalidToken.Error())
	}

	return claims, nil
}

func GetUserDataFromToken(ctx context.Context) requests.UserAuth {
	if auth, ok := ctx.Value(Token).(requests.UserAuth); ok {
		return auth
	}
	return requests.UserAuth{}
}

func GetUserRoleFromToken(ctx context.Context) (res constants.UserRoleMap) {
	res = make(map[constants.UserRole]bool)
	if auth, ok := ctx.Value(Token).(requests.UserAuth); ok {
		if auth.Role != "" {
			res[auth.Role] = true
		}
	}
	return res
}
