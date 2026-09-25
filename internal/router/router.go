package router

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/handlers"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/activity"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/auth"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/comment"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/post"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/saves"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/upload"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/modules/user"
	"github.com/gin-gonic/gin"
)

func registerModuleRoutes(group *gin.RouterGroup, h *handlers.Handlers) {
	auth.AuthRoutes(group, h.Auth)
	user.UserRoutes(group, h.User)
	post.PostRoutes(group, h.Post)
	comment.CommentRoutes(group, h.Comment)
	activity.ActivityRoutes(group, h.Activity)
	upload.UploadRoutes(group, h.Upload)
	saves.SavesRoutes(group, h.Saves)
}

func SetupRoutes(r *gin.Engine, h *handlers.Handlers) {
	// Base URL Group: /api/v1
	apiV1 := r.Group("/api/v1")
	registerModuleRoutes(apiV1, h)

	// Fallback alias root group for backward compatibility
	rootGroup := r.Group("")
	registerModuleRoutes(rootGroup, h)
}
