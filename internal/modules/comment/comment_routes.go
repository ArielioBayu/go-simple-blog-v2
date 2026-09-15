package comment

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.Engine, h *CommentHandler) {
	group := r.Group("/posts")
	group.Use(middleware.AuthMiddlewareToken())
	group.POST("/create-comment/:postId", h.CreateComment)
	group.GET("/comment-count/:postId", h.CountComments)
}
