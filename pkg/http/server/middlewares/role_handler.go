package middlewares

import (
	"backend-ta/app/constants"
	"backend-ta/pkg/authentication"
	"backend-ta/pkg/errors"
	"backend-ta/pkg/http/server/http_response"

	"github.com/gin-gonic/gin"
)

func RoleHandler(roles ...constants.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(roles) == 0 {
			c.Next()
			return
		}

		userRoles := authentication.GetUserRoleFromToken(c)
		for _, role := range roles {
			if userRoles[role] {
				c.Next()
				return
			}
		}

		http_response.SendError(c, errors.ForbiddenErrorToAppError())
	}
}
