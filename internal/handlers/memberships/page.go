package memberships

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Page(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "This is the first page",
	})
}
