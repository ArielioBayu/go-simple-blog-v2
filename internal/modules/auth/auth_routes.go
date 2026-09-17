package auth

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup, h *AuthHandler) {
	group := r.Group("/auth")
	group.POST("/sign-up", h.SignUp)
	group.POST("/sign-in", h.SignIn)
	group.POST("/refresh", h.Refresh)
	group.POST("/sign-out", middleware.AuthMiddlewareToken(), h.SignOut)
}
