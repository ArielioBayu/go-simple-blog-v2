package auth

import (
	"errors"
	"log"
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService AuthService
	cfg         *configs.Config
}

func NewAuthHandler(authService AuthService, cfg *configs.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		cfg:         cfg,
	}
}

func (h *AuthHandler) SignIn(c *gin.Context) {
	ctx := c.Request.Context()

	var request SignInRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	token, refreshToken, err := h.authService.SignIn(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrDataNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}
		if errors.Is(err, constants.ErrInvalidPassword) {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})
			return
		}

		log.Printf("handler sign in failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	utils.SetAccessTokenCookie(c, token)

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "Success Login",
		Data: gin.H{
			"access_token":  token,
			"refresh_token": refreshToken,
		},
	})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var request SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	err := h.authService.SignUp(ctx, request)
	if err != nil {
		if errors.Is(err, constants.ErrUsernameOrEmailAlreadyExists) {
			c.JSON(http.StatusConflict, response.MessageResponse{
				Status:  http.StatusConflict,
				Message: err.Error(),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{
		Status:  http.StatusCreated,
		Message: "success created data",
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var request RefreshTokenRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	refreshToken, err := h.authService.GetIdRefreshToken(ctx, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRefreshTokenNotFound), errors.Is(err, constants.ErrDataNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		default:
			log.Printf("handler refresh failed | error: %v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
			return
		}
	}

	token, err := h.authService.ValidateRefreshToken(ctx, refreshToken.UserId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrRefreshTokenNotFound), errors.Is(err, constants.ErrDataNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrTokenExpired):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		case errors.Is(err, constants.ErrInvalidToken):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})
			return

		default:
			log.Printf("handler refresh failed | error: %v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "internal server error",
			})
			return
		}
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success refresh token",
		Data: gin.H{
			"access_token": token,
		},
	})
}
