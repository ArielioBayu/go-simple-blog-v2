package middleware

import (
	"net/http"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

func extractToken(ctx *gin.Context) string {
	header := strings.TrimSpace(ctx.Request.Header.Get("Authorization"))
	if header != "" {
		tokenStr := strings.TrimPrefix(header, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)
		if tokenStr != "" {
			return tokenStr
		}
	}

	cookie, err := ctx.Cookie("access_token")
	if err == nil && strings.TrimSpace(cookie) != "" {
		return strings.TrimSpace(cookie)
	}

	return ""
}

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(ctx *gin.Context) {
		tokenStr := extractToken(ctx)
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: constants.ErrMissingToken.Error(),
			})
			return
		}

		id, username, err := jwt.ValidateToken(tokenStr, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: constants.ErrInvalidToken.Error(),
			})
			return
		}
		ctx.Set("id", id)
		ctx.Set("username", username)
		ctx.Next()
	}
}

func AuthMiddlewareToken() gin.HandlerFunc {
	return AuthMiddleware()
}
