package middlewares

import (
	"backend-ta/app/dto/requests"
	"backend-ta/pkg/authentication"
	internal_err "backend-ta/pkg/errors"
	"backend-ta/pkg/http/server/http_response"
	"strings"

	"github.com/gin-gonic/gin"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			http_response.SendError(c, internal_err.AuthError(authentication.AuthErrMalformedToken.Error()))
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			http_response.SendError(c, internal_err.AuthError(authentication.AuthErrInvalidToken.Error()))
			return
		}

		accessToken, err := authentication.JWTAuth.VerifyAccessToken(tokenString)
		if err != nil {
			http_response.SendError(c, err)
			return
		}

		c.Set(authentication.Token, requests.UserAuth{
			UserID:  accessToken.UserID,
			StoreID: accessToken.StoreID,
			Email:   accessToken.Email,
			Role:    accessToken.Role,
		})

		c.Next()
	}
}
