package router

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, h *handlers.Handlers) {
	auth.AuthRoutes(r, h.Auth)
	user.UserRoutes(r, h.User)
	post.PostRoutes(r, h.Post)
	comment.CommentRoutes(r, h.Comment)
	activity.ActivityRoutes(r, h.Activity)
}
