package saves

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

type SavesHandler struct {
	savesService SavesService
	cfg          *configs.Config
}

func NewSavesHandler(savesService SavesService, cfg *configs.Config) *SavesHandler {
	return &SavesHandler{
		savesService: savesService,
		cfg:          cfg,
	}
}

func (h *SavesHandler) SavePost(c *gin.Context) {
	ctx := c.Request.Context()

	var request SaveRequest
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
			Message: "Invalid Post Id",
		})
		return
	}

	userId, exists := c.Get("id")
	if !exists || userId.(int) == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	err = h.savesService.SavePost(ctx, postId, userId.(int), request)
	if err != nil {
		if errors.Is(err, constants.ErrPostNotFound) {
			c.JSON(http.StatusNotFound, response.MessageResponse{
				Status:  http.StatusNotFound,
				Message: "Post Not Found",
			})
			return
		}

		log.Printf("SavePost failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Status:  http.StatusOK,
		Message: "success",
	})
}

func (h *SavesHandler) GetSavedPosts(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
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

	res, err := h.savesService.GetSavedPosts(ctx, userID, page, limit)
	if err != nil {
		log.Printf("GetSavedPosts failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     http.StatusOK,
		"message":    "success get saved posts",
		"data":       res.Data,
		"pagination": res.Pagination,
	})
}

func (h *SavesHandler) GetSavedPostIDs(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.GetInt("id")
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, response.MessageResponse{
			Status:  http.StatusUnauthorized,
			Message: "Unauthorized",
		})
		return
	}

	ids, err := h.savesService.GetSavedPostIDs(ctx, userID)
	if err != nil {
		log.Printf("GetSavedPostIDs failed | error: %v", err)
		c.JSON(http.StatusInternalServerError, response.MessageResponse{
			Status:  http.StatusInternalServerError,
			Message: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response.DataResponse{
		Status:  http.StatusOK,
		Message: "success get saved post ids",
		Data:    ids,
	})
}

