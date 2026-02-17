package memberships

import (
	"net/http"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/memberships"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	// binding json
	var request memberships.SignUpRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	data, err := h.membershipsService.GetUser(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
			"data":    data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"message": "success get data",
			"data":    data,
		})
	}
}
