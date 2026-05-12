package middleware

import (
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/service"
	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(
	authService *service.AuthService,
) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			utils.ErrorResponse(
				c,
				401,
				"Authorization header missing",
			)
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {

			utils.ErrorResponse(
				c,
				401,
				"Invalid authorization header",
			)

			c.Abort()
			return
		}

		token := tokenParts[1]

		claims, err := authService.ValidateToken(token)

		if err != nil {

			utils.ErrorResponse(
				c,
				401,
				"Invalid token",
			)

			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}