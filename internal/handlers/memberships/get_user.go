package memberships

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "unauthorized",
		})
		return
	}

	data, err := h.membershipsService.GetUserById(ctx, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	if data == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success get data",
		"data": gin.H{
			"id":         data.ID,
			"username":   data.Username,
			"email":      data.Email,
			"created_at": data.CreatedAt,
		},
	})
}
