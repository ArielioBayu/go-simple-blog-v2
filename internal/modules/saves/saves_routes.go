package saves

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func SavesRoutes(r *gin.RouterGroup, h *SavesHandler) {
	postsGroup := r.Group("/posts")
	postsGroup.Use(middleware.AuthMiddlewareToken())
	{
		postsGroup.POST("/bookmarks/:postId", h.SavePost)
		postsGroup.GET("/bookmarks", h.GetSavedPosts)
		postsGroup.GET("/bookmarks/ids", h.GetSavedPostIDs)
	}
}
