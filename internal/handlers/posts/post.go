package posts

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
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

func (h *Handler) GetPostById(c *gin.Context) {
	ctx := c.Request.Context()
	postId := c.Param("postId")

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid post id",
		})
		return
	}

	data, err := h.srv.GetPostById(ctx, postIdInt)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "Post Not Found",
			})
			return
		}

		log.Printf("GetPostById Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
