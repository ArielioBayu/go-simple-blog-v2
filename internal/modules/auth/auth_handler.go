package auth

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
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
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	token, refreshToken, err := h.authService.SignIn(ctx, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	utils.SetAccessTokenCookie(c, token)

	response.Data(c, http.StatusOK, "Success Login", gin.H{
		"access_token":  token,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	ctx := c.Request.Context()

	var request SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	err := h.authService.SignUp(ctx, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success created data, please check your email for the OTP verification code")
}

func (h *AuthHandler) VerifyOTP(c *gin.Context) {
	ctx := c.Request.Context()

	var request VerifyOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	err := h.authService.VerifyOTP(ctx, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "email successfully verified, please login")
}

func (h *AuthHandler) ResendOTP(c *gin.Context) {
	ctx := c.Request.Context()

	var request ResendOTPRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	err := h.authService.ResendOTP(ctx, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "new OTP verification code has been sent to your email")
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	ctx := c.Request.Context()

	var request RefreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	refreshToken, err := h.authService.GetIdRefreshToken(ctx, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	token, err := h.authService.ValidateRefreshToken(ctx, refreshToken.UserId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success refresh token", gin.H{
		"access_token": token,
	})
}

func (h *AuthHandler) SignOut(c *gin.Context) {
	ctx := c.Request.Context()
	userId := c.GetInt("id")

	if userId > 0 {
		_ = h.authService.SignOut(ctx, userId)
	}

	utils.ClearCookie(c, "access_token")

	response.Success(c, http.StatusOK, "success logout")
}
