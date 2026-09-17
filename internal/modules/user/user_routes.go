package user

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup, h *UserHandler) {
	group := r.Group("/accounts")
	group.Use(middleware.AuthMiddlewareToken())

	group.GET("/user", h.GetUser)
	group.GET("/profile", h.GetProfile)
	group.GET("/profile/:id", h.GetProfileByID)
	group.PUT("/edit/profile", h.UpdateProfile)
}
