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
		tokenString, err := c.Cookie("access_token")
		if err != nil || tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if tokenString == "" {
			http_response.SendError(c, internal_err.AuthError(authentication.AuthErrMalformedToken.Error()))
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
