package memberships

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
}

func NewHandler(api *gin.Engine) *Handler {
	return &Handler{
		Engine: api,
	}
}

func (h *Handler) RegisterRoute() {
	routes := h.Group("/memberships")
	routes.GET("/page", h.Page)
}
