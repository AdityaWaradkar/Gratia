package handler

import (
	"strings"

	"github.com/adityawaradkar/gratia/auth_service/internal/model"
	"github.com/adityawaradkar/gratia/auth_service/internal/service"
	"github.com/adityawaradkar/gratia/auth_service/internal/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	AuthService *service.AuthService
}

func NewAuthHandler(
	authService *service.AuthService,
) *AuthHandler {

	return &AuthHandler{
		AuthService: authService,
	}
}

// Signup godoc
// @Summary Register new user
// @Description Creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.SignupRequest true "Signup Request"
// @Success 201 {object} model.MessageResponse
// @Failure 400 {object} model.ErrorResponse
// @Failure 409 {object} model.ErrorResponse
// @Router /auth/signup [post]
func (h *AuthHandler) Signup(c *gin.Context) {
	var request model.SignupRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		utils.ErrorResponse(
			c,
			400,
			"Invalid request data",
		)
		return
	}

	err = h.AuthService.Signup(
		c.Request.Context(),
		&request,
	)

	if err != nil {

		if err.Error() == "email already exists" {
			utils.ErrorResponse(
				c,
				409,
				err.Error(),
			)
			return
		}

		utils.ErrorResponse(
			c,
			500,
			"Internal server error",
		)
		return
	}

	utils.SuccessResponse(
		c,
		201,
		"User created successfully",
	)
}

// Login godoc
// @Summary Login user
// @Description Authenticates user and returns JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body model.LoginRequest true "Login Request"
// @Success 200 {object} model.AuthResponse
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var request model.LoginRequest

	err := c.ShouldBindJSON(&request)

	if err != nil {
		utils.ErrorResponse(
			c,
			400,
			"Invalid request data",
		)
		return
	}

	response, err := h.AuthService.Login(
		c.Request.Context(),
		&request,
	)

	if err != nil {
		utils.ErrorResponse(
			c,
			401,
			"Invalid credentials",
		)
		return
	}

	c.JSON(200, response)
}

// Validate godoc
// @Summary Validate JWT token
// @Description Validates JWT token and extracts claims
// @Tags Auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} model.ErrorResponse
// @Router /auth/validate [get]
func (h *AuthHandler) Validate(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		utils.ErrorResponse(
			c,
			401,
			"Authorization header missing",
		)
		return
	}

	tokenParts := strings.Split(authHeader, " ")

	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		utils.ErrorResponse(
			c,
			401,
			"Invalid authorization header",
		)
		return
	}

	token := tokenParts[1]

	claims, err := h.AuthService.ValidateToken(token)

	if err != nil {
		utils.ErrorResponse(
			c,
			401,
			"Invalid token",
		)
		return
	}

	c.JSON(200, gin.H{
		"user_id": claims.UserID,
		"role":    claims.Role,
	})
}