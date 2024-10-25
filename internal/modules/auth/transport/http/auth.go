//go:generate mockgen -source=auth.go -destination=auth_mock.go -package=transporthttp
package transporthttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"booking-event/internal/common/handler"
	commonmodel "booking-event/internal/common/model"
	"booking-event/internal/modules/auth/model"
)

type AuthHandler interface {
	LoginByEmail(ctx context.Context, email string, password string) (*model.User, *model.TokenPair, error)
}

type AuthHttpHandler struct {
	authService AuthHandler
}

func NewAuthHandler(authService AuthHandler) handler.HttpHandler {
	return &AuthHttpHandler{authService: authService}
}

func (h *AuthHttpHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.LoginByEmail)
}

func (h *AuthHttpHandler) LoginByEmail(c *gin.Context) {
	var request model.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	user, tokenPair, err := h.authService.LoginByEmail(c.Request.Context(), request.Email, request.Password)
	if errors.Is(err, model.ErrInvalidPassword) {
		c.JSON(http.StatusUnauthorized, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, commonmodel.Response{
			Success: false,
			Data:    nil,
			Message: err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, commonmodel.Response{
		Success: true,
		Message: "login success",
		Data: model.LoginResponse{
			Email:           user.Email,
			UserID:          user.ID,
			Role:            string(user.Role),
			AccessToken:     tokenPair.AccessToken,
			RefreshToken:    tokenPair.RefreshToken,
			ExpAccessToken:  tokenPair.AccessTokenExp,
			ExpRefreshToken: tokenPair.RefreshTokenExp,
		},
	})
}
