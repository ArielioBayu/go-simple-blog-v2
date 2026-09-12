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

func (h *Handler) CreatePost(c *gin.Context) {
	var request posts.PostRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.MessageResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	userId := c.GetInt("id")
	err = h.srv.CreatePost(c, userId, request)
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

func (h *Handler) GetAllPost(c *gin.Context) {
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
	postResp, err := h.srv.GetAllPost(ctx, pageSize, pageIndex, userID)
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

func (h *Handler) GetPostById(c *gin.Context) {
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
	data, err := h.srv.GetPostById(ctx, postIdInt, userID)
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
