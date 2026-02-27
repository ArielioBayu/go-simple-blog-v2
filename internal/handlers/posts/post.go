package posts

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreatePost(c *gin.Context) {
	var request posts.PostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	userId := c.GetInt("id")
	err = h.srv.CreatePost(c, userId, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Success Create Post",
	})
}

func (h *Handler) GetAllPost(c *gin.Context) {
	ctx := c.Request.Context()
	pageIndexStr := c.Query("pageIndex")
	pageSizeStr := c.Query("pageSize")

	pageIndex, err := strconv.Atoi(pageIndexStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page index",
		})
		return
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid page size",
		})
		return
	}

	response, err := h.srv.GetAllPost(ctx, pageSize, pageIndex)
	if err != nil {
		log.Printf("GetAllPost Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": response,
	})
}
