package utils

import (
	"github.com/adityawaradkar/gratia/auth_service/internal/model"
	"github.com/gin-gonic/gin"
)

func SuccessResponse(
	c *gin.Context,
	statusCode int,
	message string,
) {

	c.JSON(
		statusCode,
		model.MessageResponse{
			Message: message,
		},
	)
}

func ErrorResponse(
	c *gin.Context,
	statusCode int,
	message string,
) {

	c.JSON(
		statusCode,
		model.ErrorResponse{
			Error: message,
		},
	)
}