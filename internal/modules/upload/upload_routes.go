package upload

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	"github.com/gin-gonic/gin"
)

func UploadRoutes(r *gin.RouterGroup, h *UploadHandler) {
	group := r.Group("/upload")
	group.Use(middleware.AuthMiddlewareToken())

	group.POST("", h.UploadFile)
	group.GET("/my-media", h.GetMyUploads)
}
