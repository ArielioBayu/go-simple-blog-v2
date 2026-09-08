package middleware

import (
	"net/http"
	"strings"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(ctx *gin.Context) {
		header := ctx.Request.Header.Get("Authorization")
		header = strings.TrimSpace(header)
		if header == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": constants.ErrMissingToken.Error(),
			})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		tokenStr = strings.TrimSpace(tokenStr)
		if tokenStr == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": constants.ErrMissingToken.Error(),
			})
			return
		}

		id, username, err := jwt.ValidateToken(tokenStr, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": constants.ErrInvalidToken.Error(),
			})
			return
		}
		ctx.Set("id", id)
		ctx.Set("username", username)
		ctx.Next()
	}
}

func AuthMiddlewareToken() gin.HandlerFunc {
	secretKey := configs.Get().Service.SecretKey
	return func(ctx *gin.Context) {
		accessToken, err := ctx.Cookie("access_token")
		if err != nil || strings.TrimSpace(accessToken) == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": constants.ErrMissingToken.Error(),
			})
			return
		}

		id, username, err := jwt.ValidateToken(accessToken, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": constants.ErrInvalidToken.Error(),
			})
			return
		}
		ctx.Set("id", id)
		ctx.Set("username", username)
		ctx.Next()
	}
}
