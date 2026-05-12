package middleware

import (
	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/gin-gonic/gin"
)

func RBACMiddleware(
	allowedRoles ...string,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		roleValue, exists := c.Get("role")

		if !exists {

			utils.ErrorResponse(
				c,
				401,
				"Unauthorized",
			)

			c.Abort()
			return
		}

		userRole, ok := roleValue.(string)

		if !ok {

			utils.ErrorResponse(
				c,
				401,
				"Unauthorized",
			)

			c.Abort()
			return
		}

		for _, role := range allowedRoles {

			if role == userRole {
				c.Next()
				return
			}
		}

		utils.ErrorResponse(
			c,
			403,
			"Forbidden",
		)

		c.Abort()
	}
}