package memberships

import (
	service "github.com/ArielioBayu/go-simple-blog-v2/internal/service/memberships"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	*gin.Engine
	membershipsService service.MembershipsService
}

func NewHandler(api *gin.Engine, membershipsSrv service.MembershipsService) *Handler {
	return &Handler{
		Engine:             api,
		membershipsService: membershipsSrv,
	}
}

func (h *Handler) RegisterRoute() {
	routes := h.Group("/memberships")
	routes.POST("/sign-up", h.SignUp)
	routes.POST("/sign-in", h.SignIn)
	routes.GET("/get-user", h.GetUser)
	routes.POST("/refresh", h.Refresh)
}
