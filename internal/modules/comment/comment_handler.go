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
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

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

func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid post id",
		})
		return
	}

	page := 1
	if pStr := c.Query("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	} else if pStr := c.Query("pageIndex"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	} else if lStr := c.Query("pageSize"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	res, err := h.commentService.GetCommentsByPostID(c.Request.Context(), postId, page, limit)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}

		log.Printf("GetCommentsByPostID error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get comments",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *CommentHandler) CountComments(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid post id",
		})
		return
	}

	count, err := h.commentService.CountComments(c.Request.Context(), postId)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}

		log.Printf("CountComments error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get comment count",
		Data: gin.H{
			"post_id":       postId,
			"comment_count": count,
		},
	})
}

func (h *CommentHandler) CreateReply(c *gin.Context) {
	var request ReplyRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	err = h.commentService.CreateReply(c.Request.Context(), commentId, userId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrCommentNotFound):
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
			log.Printf("CreateReply error: %v", err)
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{
		Status:  http.StatusCreated,
		Message: "success create reply",
	})
}

func (h *CommentHandler) GetReplies(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
		})
		return
	}

	page := 1
	if pStr := c.Query("page"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	} else if pStr := c.Query("pageIndex"); pStr != "" {
		if p, err := strconv.Atoi(pStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 10
	if lStr := c.Query("limit"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	} else if lStr := c.Query("pageSize"); lStr != "" {
		if l, err := strconv.Atoi(lStr); err == nil && l > 0 {
			limit = l
		}
	}

	res, err := h.commentService.GetRepliesByCommentID(c.Request.Context(), commentId, page, limit)
	if err != nil {
		if errors.Is(err, constants.ErrCommentNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
			return
		}

		log.Printf("GetReplies error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get comment replies",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	var request UpdateCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: constants.ErrUnauthorized.Error(),
		})
		return
	}

	err = h.commentService.UpdateComment(c.Request.Context(), commentId, userId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, constants.ErrForbidden):
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden: you are not authorized to update this comment",
			})
		default:
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success update comment",
	})
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid comment id",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: constants.ErrUnauthorized.Error(),
		})
		return
	}

	err = h.commentService.DeleteComment(c.Request.Context(), commentId, userId)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrCommentNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, constants.ErrForbidden):
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden: you are not authorized to delete this comment",
			})
		default:
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success delete comment",
	})
}

func (h *CommentHandler) UpdateReply(c *gin.Context) {
	var request UpdateReplyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid reply id",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: constants.ErrUnauthorized.Error(),
		})
		return
	}

	err = h.commentService.UpdateReply(c.Request.Context(), replyId, userId, request)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrReplyNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, constants.ErrForbidden):
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden: you are not authorized to update this reply",
			})
		default:
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success update reply",
	})
}

func (h *CommentHandler) DeleteReply(c *gin.Context) {
	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "Invalid reply id",
		})
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: constants.ErrUnauthorized.Error(),
		})
		return
	}

	err = h.commentService.DeleteReply(c.Request.Context(), replyId, userId)
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrReplyNotFound):
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
			})
		case errors.Is(err, constants.ErrForbidden):
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: "forbidden: you are not authorized to delete this reply",
			})
		default:
			c.JSON(http.StatusInternalServerError, response.MessageResponse{
				Status:  http.StatusInternalServerError,
				Message: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success delete reply",
	})
}
