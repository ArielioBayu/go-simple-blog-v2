package user

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, h *UserHandler) {
	group := r.Group("/memberships")
	group.GET("/get-user", middleware.AuthMiddlewareToken(), h.GetUser)
}
