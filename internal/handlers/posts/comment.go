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

func (h *Handler) CreateComment(c *gin.Context) {
	var request posts.CommentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	//	GetPostId dari param
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post id",
		})
		return
	}

	//	GetUserId dari middleware menggunakan context
	userId, exists := c.Get("id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	err = h.srv.CreateComment(c, postId, userId.(int), request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrPostNotFound):
			c.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
		case errors.Is(err, constants.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": err.Error(),
			})
		default:
			//	Log agar detail error internal server tidak sampai ke client
			log.Printf("Create comment error :%v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "success create comment",
	})

}
