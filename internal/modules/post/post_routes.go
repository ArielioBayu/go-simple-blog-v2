package post

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func PostRoutes(r *gin.Engine, h *PostHandler) {
	group := r.Group("/posts")
	group.Use(middleware.AuthMiddlewareToken())

	group.POST("/create-post", h.CreatePost)
	group.GET("/get-all-post", h.GetAllPost)
	group.GET("/get-post-by-id/:postId", h.GetPostById)
}
