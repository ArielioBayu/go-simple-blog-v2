package comment

import (
	"net/http"
	"strconv"

	"github.com/ArielioBayu/go-simple-blog-v2/internal/configs"
	"github.com/ArielioBayu/go-simple-blog-v2/internal/constants"
	"github.com/ArielioBayu/go-simple-blog-v2/pkg/apperror"
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
		response.Error(c, err)
		return
	}

	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid post id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.CreateComment(c.Request.Context(), postId, userId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success create comment")
}

func (h *CommentHandler) GetCommentsByPostID(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid post id", nil))
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
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get comments", res.Data, res.Pagination)
}

func (h *CommentHandler) CountComments(c *gin.Context) {
	postId, err := strconv.Atoi(c.Param("postId"))
	if err != nil || postId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid post id", nil))
		return
	}

	count, err := h.commentService.CountComments(c.Request.Context(), postId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Data(c, http.StatusOK, "success get comment count", gin.H{
		"post_id":       postId,
		"comment_count": count,
	})
}

func (h *CommentHandler) CreateReply(c *gin.Context) {
	var request ReplyRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response.Error(c, err)
		return
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil || commentId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid comment id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.CreateReply(c.Request.Context(), commentId, userId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "success create reply")
}

func (h *CommentHandler) GetReplies(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil || commentId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid comment id", nil))
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
		response.Error(c, err)
		return
	}

	response.Paginated(c, http.StatusOK, "success get comment replies", res.Data, res.Pagination)
}

func (h *CommentHandler) UpdateComment(c *gin.Context) {
	var request UpdateCommentRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil || commentId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid comment id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.UpdateComment(c.Request.Context(), commentId, userId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success update comment")
}

func (h *CommentHandler) DeleteComment(c *gin.Context) {
	commentId, err := strconv.Atoi(c.Param("commentId"))
	if err != nil || commentId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid comment id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.DeleteComment(c.Request.Context(), commentId, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success delete comment")
}

func (h *CommentHandler) UpdateReply(c *gin.Context) {
	var request UpdateReplyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, err)
		return
	}

	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil || replyId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid reply id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.UpdateReply(c.Request.Context(), replyId, userId, request)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success update reply")
}

func (h *CommentHandler) DeleteReply(c *gin.Context) {
	replyId, err := strconv.Atoi(c.Param("replyId"))
	if err != nil || replyId <= 0 {
		response.Error(c, apperror.NewBadRequest("Invalid reply id", nil))
		return
	}

	userId := c.GetInt("id")
	if userId == 0 {
		response.Error(c, constants.ErrUnauthorized)
		return
	}

	err = h.commentService.DeleteReply(c.Request.Context(), replyId, userId)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, http.StatusOK, "success delete reply")
}
