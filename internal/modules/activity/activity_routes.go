package activity

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func ActivityRoutes(r *gin.Engine, h *ActivityHandler) {
	group := r.Group("/posts")
	group.Use(middleware.AuthMiddlewareToken())
	group.POST("/user-activity/:postId", h.InsertUpdateActivities)
}
