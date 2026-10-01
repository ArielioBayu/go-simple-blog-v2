package comment

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func CommentRoutes(r *gin.RouterGroup, h *CommentHandler) {
	// Post Comments
	commentGroup := r.Group("/posts/comments")
	commentGroup.Use(middleware.AuthMiddlewareToken())
	{
		commentGroup.POST("/:postId", h.CreateComment)
		commentGroup.GET("/:postId", h.GetCommentsByPostID)
		commentGroup.GET("/count/:postId", h.CountComments)
		commentGroup.PUT("/:commentId", h.UpdateComment)
		commentGroup.DELETE("/:commentId", h.DeleteComment)
	}

	// Comment Replies
	repliesGroup := r.Group("/comments/replies")
	repliesGroup.Use(middleware.AuthMiddlewareToken())
	{
		repliesGroup.POST("/:commentId", h.CreateReply)
		repliesGroup.GET("/:commentId", h.GetReplies)
		repliesGroup.PUT("/:replyId", h.UpdateReply)
		repliesGroup.DELETE("/:replyId", h.DeleteReply)
	}
}
