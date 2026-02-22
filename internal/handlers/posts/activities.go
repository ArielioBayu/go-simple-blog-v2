package posts

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/gin-gonic/gin"
)

func (h *Handler) InsertUpdateActivities(c *gin.Context) {
	ctx := c.Request.Context()

	var request posts.ActivityRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"messsage": err.Error(),
		})
		return
	}

	// Get postId from param
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid Post Id",
		})
		return
	}

	// Get userId from middleware
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	err = h.srv.InsertUpdateActivities(ctx, postId, userId.(int), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
