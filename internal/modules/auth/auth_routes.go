package auth

import "github.com/gin-gonic/gin"

func AuthRoutes(r *gin.Engine, h *AuthHandler) {
	group := r.Group("/memberships")
	group.POST("/sign-up", h.SignUp)
	group.POST("/sign-in", h.SignIn)
	group.POST("/refresh", h.Refresh)
}
