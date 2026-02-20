package posts

import (
	"github.com/ArielioBayu/go-simple-blog-v2/internal/middleware"
	service "github.com/ArielioBayu/go-simple-blog-v2/internal/service/posts"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
	srv service.PostsService
}

func NewHandler(api *gin.Engine, service service.PostsService) *Handler {
	return &Handler{
		Engine: api,
		srv:    service,
	}
}

func (h *Handler) RegisterRoute() {
	routes := h.Group("/posts")
	routes.Use(middleware.AuthMiddleware())

	routes.POST("/create-post", h.CreatePost)
	routes.POST("/create-comment/:postId", h.CreateComment)
}
