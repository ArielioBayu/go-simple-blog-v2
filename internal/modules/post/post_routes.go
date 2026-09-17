package post

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PostRoutes(r *gin.RouterGroup, h *PostHandler) {
	group := r.Group("/posts")
	group.Use(middleware.AuthMiddlewareToken())

	group.POST("", h.CreatePost)
	group.GET("", h.GetAllPost)
	group.GET("/:postId", h.GetPostById)
	group.DELETE("/:postId", h.DeletePost)
	group.GET("/user/:userId", h.GetPostsByUserID)
}
