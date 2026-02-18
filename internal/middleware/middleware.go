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
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, constants.ErrMissingToken)
			return
		}

		id, username, err := jwt.ValidateToken(header, secretKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, constants.ErrInvalidToken)
			return
		}
		ctx.Set("id", id)
		ctx.Set("username", username)
		ctx.Next()
	}
}
