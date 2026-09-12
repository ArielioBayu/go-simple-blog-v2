package posts

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/model/posts"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) CreateComment(c *gin.Context) {
	var request posts.CommentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	//	GetPostId dari param
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid post id",
		})
		return
	}

	//	GetUserId dari middleware menggunakan context
	userId := c.GetInt("id")

	err = h.srv.CreateComment(c, postId, userId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrPostNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})

		case errors.Is(err, constants.ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, response.MessageResponse{
				Status:  http.StatusUnauthorized,
				Message: err.Error(),
			})

		default:
			//	Log agar detail error internal server tidak sampai ke client
			log.Printf("Create comment error :%v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{
		Status:  http.StatusCreated,
		Message: "success create comment",
	})

}
