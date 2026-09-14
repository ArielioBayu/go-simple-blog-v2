package comment

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/response"
	"github.com/gin-gonic/gin"
)

type CommentHandler struct {
	commentService CommentService
	cfg            *configs.Config
}

func NewCommentHandler(commentService CommentService, cfg *configs.Config) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		cfg:            cfg,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
	var request CommentRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid post id",
		})
		return
	}

	userId := c.GetInt("id")

	err = h.commentService.CreateComment(c.Request.Context(), postId, userId, request)
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
			log.Printf("Create comment error: %v", err)
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
