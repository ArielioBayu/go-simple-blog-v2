package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

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

func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		response.Error(c, fmt.Errorf("failed to generate oauth state: %w", err))
		return
	}
	state := hex.EncodeToString(b)

	// Save state in HTTP-only short-lived cookie for CSRF protection
	utils.SetCookie(c, "oauth_state", state, 5*time.Minute)

	url := h.authService.GetGoogleAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	stateQuery := c.Query("state")
	stateCookie, err := c.Cookie("oauth_state")

	// Clear state cookie immediately
	utils.ClearCookie(c, "oauth_state")

	frontendURL := h.cfg.Service.FrontendURL
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	if err != nil || stateQuery == "" || stateQuery != stateCookie {
		redirectURL := fmt.Sprintf("%s/login?error=%s", frontendURL, url.QueryEscape("invalid or expired oauth state"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	code := c.Query("code")
	if code == "" {
		redirectURL := fmt.Sprintf("%s/login?error=%s", frontendURL, url.QueryEscape("missing oauth authorization code"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	token, refreshToken, err := h.authService.HandleGoogleCallback(c.Request.Context(), code)
	if err != nil {
		log.Printf("[GoogleCallback] error: %v", err)
		redirectURL := fmt.Sprintf("%s/login?error=%s", frontendURL, url.QueryEscape("failed to authenticate with google"))
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	// Set HttpOnly cookies
	utils.SetAccessTokenCookie(c, token)
	utils.SetCookie(c, "refresh_token", refreshToken, 7*24*time.Hour)

	// Redirect to frontend auth callback URL with tokens
	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&refresh_token=%s", frontendURL, url.QueryEscape(token), url.QueryEscape(refreshToken))
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

