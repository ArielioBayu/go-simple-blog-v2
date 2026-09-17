package post

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

type PostHandler struct {
	postService PostService
	cfg         *configs.Config
}

func NewPostHandler(postService PostService, cfg *configs.Config) *PostHandler {
	return &PostHandler{
		postService: postService,
		cfg:         cfg,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	var request PostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	userId := c.GetInt("id")
	err = h.postService.CreatePost(c.Request.Context(), userId, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response.MessageResponse{
		Status:  http.StatusCreated,
		Message: "Success Create Post",
	})
}

func (h *PostHandler) GetAllPost(c *gin.Context) {
	ctx := c.Request.Context()

	pageIndex := 1
	if pageIndexStr := c.Query("pageIndex"); pageIndexStr != "" {
		p, err := strconv.Atoi(pageIndexStr)
		if err != nil || p <= 0 {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid page index",
			})
			return
		}
		pageIndex = p
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		s, err := strconv.Atoi(pageSizeStr)
		if err != nil || s <= 0 {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid page size",
			})
			return
		}
		pageSize = s
	}

	userID := c.GetInt("id")
	postResp, err := h.postService.GetAllPost(ctx, pageSize, pageIndex, userID)
	if err != nil {
		log.Printf("GetAllPost Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.PaginationResponse{
		Status:  http.StatusOK,
		Message: "success get all post",
		Pagination: &response.Pagination{
			Limit:  postResp.Pagination.Limit,
			Offset: postResp.Pagination.Offset,
		},
		Data: postResp.Data,
	})
}

func (h *PostHandler) GetPostById(c *gin.Context) {
	ctx := c.Request.Context()
	postId := c.Param("postId")

	postIdInt, err := strconv.Atoi(postId)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid post id",
		})
		return
	}

	userID := c.GetInt("id")
	data, err := h.postService.GetPostById(ctx, postIdInt, userID)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "Post Not Found",
			})
			return
		}

		log.Printf("GetPostById Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get post",
		Data:    data,
	})
}

func (h *PostHandler) GetPostsByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	userIdParam := c.Param("userId")
	if userIdParam == "" {
		userIdParam = c.Param("user_id")
	}

	targetUserID, err := strconv.Atoi(userIdParam)
	if err != nil || targetUserID <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid user id",
		})
		return
	}

	pageIndex := 1
	if pageIndexStr := c.Query("pageIndex"); pageIndexStr != "" {
		p, err := strconv.Atoi(pageIndexStr)
		if err != nil || p <= 0 {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid page index",
			})
			return
		}
		pageIndex = p
	} else if pageStr := c.Query("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			pageIndex = p
		}
	}

	pageSize := 10
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		s, err := strconv.Atoi(pageSizeStr)
		if err != nil || s <= 0 {
			c.JSON(http.StatusBadRequest, response.MessageResponse{
				Status:  http.StatusBadRequest,
				Message: "Invalid page size",
			})
			return
		}
		pageSize = s
	} else if limitStr := c.Query("limit"); limitStr != "" {
		s, err := strconv.Atoi(limitStr)
		if err == nil && s > 0 {
			pageSize = s
		}
	}

	currentUserID := c.GetInt("id")
	postResp, err := h.postService.GetPostsByUserID(ctx, targetUserID, currentUserID, pageSize, pageIndex)
	if err != nil {
		log.Printf("GetPostsByUserID Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.PaginationResponse{
		Status:  http.StatusOK,
		Message: "success get posts by user id",
		Pagination: &response.Pagination{
			Limit:  postResp.Pagination.Limit,
			Offset: postResp.Pagination.Offset,
		},
		Data: postResp.Data,
	})
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	ctx := c.Request.Context()

	postId := c.Param("postId")
	postIdInt, err := strconv.Atoi(postId)
	if err != nil || postIdInt <= 0 {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid post id",
		})
		return
	}

	userId := c.GetInt("id")
	err = h.postService.DeletePost(ctx, postIdInt, userId)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "post not found",
			})
			return
		}
		if errors.Is(err, constants.ErrForbidden) {
			c.JSON(http.StatusForbidden, response.MessageResponse{
				Status:  http.StatusForbidden,
				Message: "you are not allowed to delete this post",
			})
			return
		}

		log.Printf("DeletePost Failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success delete post",
	})
}


